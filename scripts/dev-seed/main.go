package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type appConfig struct {
	DB struct {
		Host        string        `koanf:"host"`
		Port        int           `koanf:"port"`
		User        string        `koanf:"user"`
		Password    string        `koanf:"password"`
		Database    string        `koanf:"database"`
		SSLMode     string        `koanf:"ssl_mode"`
		Params      string        `koanf:"params"`
		MaxOpen     int           `koanf:"max_open"`
		MaxIdle     int           `koanf:"max_idle"`
		MaxLifetime time.Duration `koanf:"max_lifetime"`
	} `koanf:"db"`
}

type seedAgent struct {
	Email     string
	FirstName string
	LastName  string
	Roles     []string
	TeamName  string
}

type seedRole struct {
	Name        string
	Description string
	Permissions []string
}

type seedInbox struct {
	Name    string
	From    string
	Channel string
}

type seedContact struct {
	Key              string
	Email            string
	FirstName        string
	LastName         string
	Phone            string
	CustomAttributes map[string]any
}

type seedMessage struct {
	Key         string
	Type        string
	SenderType  string
	SenderEmail string
	Content     string
	Offset      time.Duration
}

type seedConversation struct {
	Key              string
	Subject          string
	ContactKey       string
	InboxName        string
	Status           string
	Priority         string
	TeamName         string
	AssigneeEmail    string
	TagNames         []string
	CreatedAgo       time.Duration
	CustomAttributes map[string]any
	Messages         []seedMessage
}

func main() {
	var (
		configPath = flag.String("config", "config.toml", "caminho do config.toml")
		password   = flag.String("password", "CanalGov@123", "senha para os usuarios de desenvolvimento")
		count      = flag.Int("count", 50, "quantidade aproximada de conversas a gerar")
	)
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("erro carregando config: %v", err)
	}

	db, err := sqlx.Connect("postgres", makeDSN(cfg))
	if err != nil {
		log.Fatalf("erro conectando ao banco: %v", err)
	}
	defer db.Close()

	if err := seed(db, *password, *count); err != nil {
		log.Fatalf("erro gerando dados de desenvolvimento: %v", err)
	}

	fmt.Println("\n==================================================")
	fmt.Println("SEED DE PREFEITURA CONCLUÍDO COM SUCESSO")
	fmt.Println("==================================================")
	fmt.Printf("Senha padrão: %s\n", *password)
	fmt.Println("\nLogins principais:")
	fmt.Println("- gestor@prefeitura.gov.br (Admin)")
	fmt.Println("- saude@prefeitura.gov.br (Saúde)")
	fmt.Println("- obras@prefeitura.gov.br (Obras)")
	fmt.Println("- fazenda@prefeitura.gov.br (Fazenda)")
	fmt.Println("==================================================\n")
}

func loadConfig(path string) (*appConfig, error) {
	ko := koanf.New(".")
	if err := ko.Load(file.Provider(path), toml.Parser()); err != nil {
		return nil, err
	}
	var cfg appConfig
	if err := ko.Unmarshal("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func makeDSN(cfg *appConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s %s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Database, cfg.DB.SSLMode, cfg.DB.Params)
}

func seed(db *sqlx.DB, password string, count int) error {
	// 1. Roles
	roles := []seedRole{
		{Name: "Gestor Municipal", Description: "Acesso total à prefeitura", Permissions: []string{"*"}},
		{Name: "Atendente de Secretaria", Description: "Atendimento específico de setor", Permissions: []string{"conversations:read", "messages:write"}},
	}
	for _, r := range roles {
		ensureRole(db, r)
	}

	// 2. Teams
	teamIDs := make(map[string]int64)
	teams := []struct{ Name, Emoji string }{
		{"Saúde", "🏥"},
		{"Obras", "🚧"},
		{"Fazenda", "💰"},
		{"Educação", "📚"},
		{"Social", "🤝"},
		{"Ouvidoria", "📣"},
	}
	for _, t := range teams {
		id, _ := ensureTeam(db, t.Name, t.Emoji)
		teamIDs[t.Name] = id
	}

	// 3. Inboxes
	inboxIDs := make(map[string]int64)
	inboxNames := []string{}
	inboxes := []seedInbox{
		{Name: "WhatsApp Cidadão", From: "5511999998888", Channel: "whatsapp"},
		{Name: "Telegram Prefeitura", From: "@PrefeituraBot", Channel: "telegram"},
		{Name: "E-mail Institucional", From: "contato@prefeitura.gov.br", Channel: "email"},
	}
	for _, ib := range inboxes {
		id, _ := ensureInboxFull(db, ib)
		inboxIDs[ib.Name] = id
		inboxNames = append(inboxNames, ib.Name)
	}

	// 4. Agents
	agentIDs := make(map[string]int64)
	agents := []seedAgent{
		{Email: "gestor@prefeitura.gov.br", FirstName: "Ricardo", LastName: "Gestor", Roles: []string{"Admin", "Gestor Municipal"}, TeamName: "Ouvidoria"},
		{Email: "saude@prefeitura.gov.br", FirstName: "Dra. Helena", LastName: "Silva", Roles: []string{"Atendente de Secretaria"}, TeamName: "Saúde"},
		{Email: "obras@prefeitura.gov.br", FirstName: "Eng. Marcos", LastName: "Oliveira", Roles: []string{"Atendente de Secretaria"}, TeamName: "Obras"},
		{Email: "fazenda@prefeitura.gov.br", FirstName: "Sérgio", LastName: "Fiscal", Roles: []string{"Atendente de Secretaria"}, TeamName: "Fazenda"},
	}
	for _, a := range agents {
		id, _ := ensureAgent(db, a, password)
		agentIDs[a.Email] = id
		ensureTeamMember(db, teamIDs[a.TeamName], id, "⭐")
	}

	// 5. Contacts (Citizens)
	citizens := []struct{ First, Last string }{
		{"João", "Silva"}, {"Maria", "Santos"}, {"José", "Oliveira"}, {"Ana", "Pereira"}, {"Carlos", "Ferreira"},
		{"Paulo", "Almeida"}, {"Lucas", "Costa"}, {"Carla", "Gomes"}, {"Luiz", "Rocha"}, {"Marcos", "Ribeiro"},
		{"Aline", "Barros"}, {"Diego", "Mendes"}, {"Patrícia", "Teixeira"}, {"Roberto", "Nunes"}, {"Fernanda", "Lima"},
	}
	contactIDs := make(map[string]int64)
	contactKeys := []string{}
	for i, c := range citizens {
		key := fmt.Sprintf("citizen_%d", i)
		contact := seedContact{
			Key:       key,
			Email:     fmt.Sprintf("%s.%s@exemplo.com", strings.ToLower(c.First), strings.ToLower(c.Last)),
			FirstName: c.First,
			LastName:  c.Last,
			Phone:     fmt.Sprintf("55119%08d", rand.Intn(99999999)),
			CustomAttributes: map[string]any{
				"cartao_sus": fmt.Sprintf("%d", 10000000+rand.Intn(90000000)),
				"nis":        fmt.Sprintf("%d", 10000000000+rand.Intn(90000000000)),
			},
		}
		
		var lastID int64
		for _, ibID := range inboxIDs {
			id, _, _ := ensureContact(db, ibID, contact)
			lastID = id
		}
		contactIDs[key] = lastID
		contactKeys = append(contactKeys, key)
	}

	// 6. Template Categories
	fmt.Println("Gerando categorias de modelos...")
	seedCats := []struct{ Name, Desc, TeamName string }{
		{"Protocolos de Saúde", "Modelos para agendamentos e triagem", "Saúde"},
		{"Vistorias e Obras", "Modelos para fiscalização de campo", "Obras"},
		{"Atendimento Geral", "Modelos universais de nota", ""},
	}
	catIDs := make(map[string]int64)
	for _, sc := range seedCats {
		var id int64
		db.Get(&id, `INSERT INTO template_categories (name, description) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id`, sc.Name, sc.Desc)
		if id == 0 { db.Get(&id, `SELECT id FROM template_categories WHERE name = $1`, sc.Name) }
		catIDs[sc.Name] = id
		if sc.TeamName != "" {
			db.Exec(`INSERT INTO template_category_teams (category_id, team_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, teamIDs[sc.TeamName])
		}
	}

	// 7. Note Templates
	fmt.Println("Gerando modelos de notas...")
	noteTpls := []struct{ Name, Body, CatName string }{
		{"Triagem Concluída", "Cidadão passou pela triagem inicial. Encaminhado para especialidade.", "Protocolos de Saúde"},
		{"Buraco Identificado", "Constatado buraco em via pública. Necessário operação tapa-buraco.", "Vistorias e Obras"},
		{"Contato Telefônico", "Realizado contato com o cidadão. Assunto pendente de documentação.", "Atendimento Geral"},
	}
	for _, nt := range noteTpls {
		db.Exec(`INSERT INTO templates (name, body, type, category_id) VALUES ($1, $2, 'note', $3) ON CONFLICT DO NOTHING`, nt.Name, nt.Body, catIDs[nt.CatName])
	}

	// 8. Conversations
	fmt.Println("Gerando conversas...")
	subjects := map[string][]string{
		"Saúde":    {"Agendamento de consulta", "Falta de medicamento", "Vacinação infantil", "Resultado de exame"},
		"Obras":    {"Buraco na rua", "Iluminação pública queimada", "Calçada danificada", "Limpeza de bueiro"},
		"Fazenda":  {"Dúvida sobre IPTU", "Parcelamento de dívida ativa", "Nota fiscal de serviço", "Alvará de funcionamento"},
		"Educação": {"Vaga em creche", "Material escolar", "Uniforme", "Transferência de escola"},
		"Social":   {"Auxílio Brasil", "Cesta básica", "Atendimento CRAS", "Cadastro Único"},
	}
	statusList := []string{"Open", "Open", "Resolved", "Closed"}
	priorityList := []string{"Low", "Medium", "High", "Urgent"}
	statusIDs, _ := loadLookup(db, "conversation_statuses")
	priorityIDs, _ := loadLookup(db, "conversation_priorities")

	for i := 0; i < count; i++ {
		teamName := teams[rand.Intn(len(teams))].Name
		subjs := subjects[teamName]
		if len(subjs) == 0 { subjs = []string{"Assunto Geral"} }
		contactKey := contactKeys[rand.Intn(len(contactKeys))]
		inboxName := inboxNames[rand.Intn(len(inboxNames))]
		
		conv := seedConversation{
			Key:           fmt.Sprintf("conv_%d", i),
			Subject:       fmt.Sprintf("%s - %d", subjs[rand.Intn(len(subjs))], 2026000+i),
			ContactKey:    contactKey,
			InboxName:     inboxName,
			Status:        statusList[rand.Intn(len(statusList))],
			Priority:      priorityList[rand.Intn(len(priorityList))],
			TeamName:      teamName,
			AssigneeEmail: agents[rand.Intn(len(agents))].Email,
			CreatedAgo:    time.Duration(rand.Intn(500)) * time.Hour,
			Messages: []seedMessage{
				{Key: "msg1", Type: "incoming", SenderType: "contact", Content: "Olá, gostaria de registrar uma solicitação.", Offset: 0},
				{Key: "msg2", Type: "outgoing", SenderType: "agent", SenderEmail: "gestor@prefeitura.gov.br", Content: "Olá! Recebemos seu pedido e encaminhamos para o setor responsável.", Offset: 30 * time.Minute},
			},
		}

		var channelID int64
		db.Get(&channelID, `SELECT id FROM contact_channels WHERE contact_id = $1 AND inbox_id = $2 LIMIT 1`, contactIDs[contactKey], inboxIDs[inboxName])

		pID := priorityIDs[conv.Priority]
		if pID == 0 { pID = priorityIDs["Medium"] } // Fallback

		convID, err := ensureConversation(db, conv, contactIDs[contactKey], channelID, inboxIDs[inboxName], teamIDs[teamName], agentIDs[conv.AssigneeEmail], statusIDs[conv.Status], pID)
		if err == nil { ensureMessages(db, conv, convID, contactIDs, agentIDs) }
		
		if i % 10 == 0 { fmt.Printf("Gerando dados... %d/%d\n", i, count) }
	}
	return nil
}

func ensureInboxFull(db *sqlx.DB, ib seedInbox) (int64, error) {
	var id int64
	err := db.Get(&id, `SELECT id FROM inboxes WHERE name = $1 LIMIT 1`, ib.Name)
	if err == nil { return id, nil }
	err = db.Get(&id, `INSERT INTO inboxes (name, channel, enabled, "from", config) VALUES ($1, $2, true, $3, '{}'::jsonb) RETURNING id`, ib.Name, ib.Channel, ib.From)
	return id, err
}

func ensureRole(db *sqlx.DB, role seedRole) (int64, error) {
	var id int64
	err := db.Get(&id, `SELECT id FROM roles WHERE name = $1`, role.Name)
	if err == nil { return id, nil }
	err = db.Get(&id, `INSERT INTO roles (name, description, permissions) VALUES ($1, $2, $3) RETURNING id`, role.Name, role.Description, pq.StringArray(role.Permissions))
	return id, err
}

func ensureTeam(db *sqlx.DB, name, emoji string) (int64, error) {
	var id int64
	err := db.Get(&id, `SELECT id FROM teams WHERE name = $1`, name)
	if err == nil { return id, nil }
	err = db.Get(&id, `INSERT INTO teams (name, emoji, conversation_assignment_type, timezone) VALUES ($1, $2, 'Manual', 'America/Sao_Paulo') RETURNING id`, name, emoji)
	return id, err
}

func ensureAgent(db *sqlx.DB, agent seedAgent, password string) (int64, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	var id int64
	err := db.Get(&id, `SELECT id FROM users WHERE email = $1 AND type = 'agent'`, agent.Email)
	if err == nil { return id, nil }
	err = db.Get(&id, `INSERT INTO users (email, type, first_name, last_name, password, enabled, availability_status) VALUES ($1, 'agent', $2, $3, $4, true, 'online') RETURNING id`, agent.Email, agent.FirstName, agent.LastName, hash)
	for _, roleName := range agent.Roles {
		db.Exec(`INSERT INTO user_roles (user_id, role_id) SELECT $1, id FROM roles WHERE name = $2`, id, roleName)
	}
	return id, err
}

func ensureTeamMember(db *sqlx.DB, teamID, userID int64, emoji string) error {
	_, err := db.Exec(`INSERT INTO team_members (team_id, user_id, emoji) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, teamID, userID, emoji)
	return err
}

func ensureContact(db *sqlx.DB, inboxID int64, contact seedContact) (int64, int64, error) {
	var contactID int64
	customAttrs, _ := json.Marshal(contact.CustomAttributes)
	err := db.Get(&contactID, `SELECT id FROM users WHERE email = $1 AND type = 'contact'`, contact.Email)
	if err != nil {
		err = db.Get(&contactID, `INSERT INTO users (email, type, first_name, last_name, phone_number, custom_attributes) VALUES ($1, 'contact', $2, $3, $4, $5::jsonb) RETURNING id`, contact.Email, contact.FirstName, contact.LastName, contact.Phone, string(customAttrs))
	}
	var channelID int64
	err = db.Get(&channelID, `SELECT id FROM contact_channels WHERE contact_id = $1 AND inbox_id = $2`, contactID, inboxID)
	if err != nil {
		db.Get(&channelID, `INSERT INTO contact_channels (contact_id, inbox_id, identifier) VALUES ($1, $2, $3) RETURNING id`, contactID, inboxID, contact.Email)
	}
	return contactID, channelID, nil
}

func loadLookup(db *sqlx.DB, table string) (map[string]int64, error) {
	rows, err := db.Queryx(fmt.Sprintf(`SELECT id, name FROM %s`, table))
	if err != nil { return nil, err }
	defer rows.Close()
	out := make(map[string]int64)
	for rows.Next() {
		var id int64
		var name string
		rows.Scan(&id, &name)
		out[name] = id
	}
	return out, nil
}

func ensureConversation(db *sqlx.DB, conv seedConversation, contactID, channelID, inboxID, teamID, assigneeID, statusID, priorityID int64) (int64, error) {
	var id int64
	createdAt := time.Now().Add(-conv.CreatedAgo)
	err := db.Get(&id, `INSERT INTO conversations (created_at, updated_at, contact_id, assigned_user_id, assigned_team_id, inbox_id, contact_channel_id, status_id, priority_id, subject, meta, custom_attributes, waiting_since)
		VALUES ($1, $1, $2, $3, $4, $5, $6, $7, $8, $9, '{}'::jsonb, '{}'::jsonb, $1) RETURNING id`,
		createdAt, contactID, assigneeID, teamID, inboxID, channelID, statusID, priorityID, conv.Subject)
	return id, err
}

func ensureMessages(db *sqlx.DB, conv seedConversation, conversationID int64, contactIDs, agentIDs map[string]int64) error {
	base := time.Now().Add(-conv.CreatedAgo)
	for _, msg := range conv.Messages {
		var senderID int64
		if msg.SenderType == "contact" {
			senderID = contactIDs[conv.ContactKey]
		} else {
			senderID = agentIDs[msg.SenderEmail]
		}
		db.Exec(`INSERT INTO conversation_messages (created_at, updated_at, type, status, private, conversation_id, content_type, content, text_content, source_id, sender_id, sender_type, meta)
			VALUES ($1, $1, $2, 'sent', false, $3, 'text', $4, $4, $5, $6, $7, '{}'::jsonb)`,
			base.Add(msg.Offset), msg.Type, conversationID, msg.Content, fmt.Sprintf("seed-%s-%s", conv.Key, msg.Key), senderID, msg.SenderType)
	}
	return nil
}
