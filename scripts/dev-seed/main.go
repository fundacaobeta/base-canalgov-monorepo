package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

type seedTemplate struct {
	Name      string
	Type      string
	Subject   string
	Body      string
	TeamName  string
	IsDefault bool
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

type seedInbox struct {
	Name string
	From string
}

type seedBusinessHours struct {
	Name         string
	Description  string
	IsAlwaysOpen bool
	Hours        map[string]map[string]string
}

type seedSLA struct {
	Name              string
	Description       string
	FirstResponseTime string
	ResolutionTime    string
	NextResponseTime  string
	Notifications     []map[string]any
}

type seedCustomAttribute struct {
	Name        string
	Description string
	AppliesTo   string
	Key         string
	Values      []string
	DataType    string
	Regex       string
	RegexHint   string
}

type seedMacro struct {
	Name           string
	MessageContent string
	Visibility     string
	VisibleWhen    []string
	TeamName       string
	AgentEmail     string
	Actions        []map[string]any
}

type seedWebhook struct {
	Name     string
	URL      string
	Events   []string
	Secret   string
	IsActive bool
}

type seedView struct {
	Name       string
	Visibility string
	TeamName   string
	AgentEmail string
	Filters    []map[string]any
}

type seedAutomationRule struct {
	Name          string
	Description   string
	Type          string
	Events        []string
	Enabled       bool
	Weight        int
	ExecutionMode string
	Rules         []map[string]any
}

func main() {
	var (
		configPath = flag.String("config", "config.toml", "caminho do config.toml")
		password   = flag.String("password", "CanalGov@123", "senha para os usuarios de desenvolvimento")
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

	db.SetMaxOpenConns(cfg.DB.MaxOpen)
	db.SetMaxIdleConns(cfg.DB.MaxIdle)
	db.SetConnMaxLifetime(cfg.DB.MaxLifetime)

	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("erro abrindo transacao: %v", err)
	}
	defer tx.Rollback()

	if err := seed(tx, *password); err != nil {
		log.Fatalf("erro gerando dados de desenvolvimento: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("erro confirmando seed: %v", err)
	}

	fmt.Println("seed de desenvolvimento concluido com sucesso")
	fmt.Println("usuarios criados/atualizados com senha:", *password)
	fmt.Println("logins: ana.silva@canalgov.local, bruno.costa@canalgov.local, carla.souza@canalgov.local")
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
	if cfg.DB.SSLMode == "" {
		cfg.DB.SSLMode = "disable"
	}
	return &cfg, nil
}

func makeDSN(cfg *appConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s %s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Database,
		cfg.DB.SSLMode,
		cfg.DB.Params,
	)
}

func seed(tx *sqlx.Tx, password string) error {
	roles := []seedRole{
		{
			Name:        "Gestor de Atendimento",
			Description: "Cargo para liderancas que acompanham conversas, filas e operacao do atendimento.",
			Permissions: []string{
				"conversations:read_all",
				"conversations:read_unassigned",
				"conversations:read_assigned",
				"conversations:read_team_inbox",
				"conversations:read_team_all",
				"conversations:read",
				"conversations:update_user_assignee",
				"conversations:update_team_assignee",
				"conversations:update_priority",
				"conversations:update_status",
				"conversations:update_tags",
				"messages:read",
				"messages:write",
				"view:manage",
				"reports:manage",
				"teams:manage",
			},
		},
		{
			Name:        "Analista de Ouvidoria",
			Description: "Cargo para analistas que tratam manifestacoes, beneficios e respostas oficiais.",
			Permissions: []string{
				"conversations:read_assigned",
				"conversations:read_team_inbox",
				"conversations:read_team_all",
				"conversations:read",
				"conversations:update_user_assignee",
				"conversations:update_priority",
				"conversations:update_status",
				"conversations:update_tags",
				"messages:read",
				"messages:write",
				"contacts:read",
				"contact_notes:read",
				"contact_notes:write",
				"view:manage",
			},
		},
		{
			Name:        "Operador de Protocolo",
			Description: "Cargo para operadores que acompanham protocolos, prazos e tratativas pendentes.",
			Permissions: []string{
				"conversations:read_assigned",
				"conversations:read_team_inbox",
				"conversations:read",
				"conversations:update_priority",
				"conversations:update_status",
				"conversations:update_tags",
				"messages:read",
				"messages:write",
				"contacts:read",
				"view:manage",
			},
		},
	}
	for _, role := range roles {
		if _, err := ensureRole(tx, role); err != nil {
			return err
		}
	}

	businessHoursIDs := make(map[string]int64)
	for _, item := range []seedBusinessHours{
		{
			Name:        "Expediente padrão",
			Description: "Horário comercial padrão do atendimento administrativo.",
			Hours: map[string]map[string]string{
				"monday":    {"open": "08:00", "close": "18:00"},
				"tuesday":   {"open": "08:00", "close": "18:00"},
				"wednesday": {"open": "08:00", "close": "18:00"},
				"thursday":  {"open": "08:00", "close": "18:00"},
				"friday":    {"open": "08:00", "close": "18:00"},
			},
		},
		{
			Name:         "Plantão da ouvidoria",
			Description:  "Horário estendido para manifestações sensíveis e retornos prioritários.",
			IsAlwaysOpen: false,
			Hours: map[string]map[string]string{
				"monday":    {"open": "07:00", "close": "20:00"},
				"tuesday":   {"open": "07:00", "close": "20:00"},
				"wednesday": {"open": "07:00", "close": "20:00"},
				"thursday":  {"open": "07:00", "close": "20:00"},
				"friday":    {"open": "07:00", "close": "20:00"},
				"saturday":  {"open": "08:00", "close": "12:00"},
			},
		},
		{
			Name:         "Cobertura integral",
			Description:  "Cobertura contínua para filas com atendimento ininterrupto.",
			IsAlwaysOpen: true,
			Hours:        map[string]map[string]string{},
		},
	} {
		id, err := ensureBusinessHours(tx, item)
		if err != nil {
			return err
		}
		businessHoursIDs[item.Name] = id
	}

	slaIDs := make(map[string]int64)
	for _, item := range []seedSLA{
		{
			Name:              "SLA atendimento inicial",
			Description:       "Meta padrão para primeiro retorno e resolução do atendimento geral.",
			FirstResponseTime: "2h",
			ResolutionTime:    "24h",
			NextResponseTime:  "8h",
			Notifications: []map[string]any{
				{"type": "warning", "time_delay_type": "before", "time_delay": "30m", "metric": "first_response", "recipients": []string{"assigned_user"}},
			},
		},
		{
			Name:              "SLA ouvidoria prioritária",
			Description:       "Meta mais restrita para manifestações prioritárias da ouvidoria.",
			FirstResponseTime: "1h",
			ResolutionTime:    "12h",
			NextResponseTime:  "4h",
			Notifications: []map[string]any{
				{"type": "warning", "time_delay_type": "before", "time_delay": "15m", "metric": "resolution", "recipients": []string{"assigned_user", "team"}},
			},
		},
		{
			Name:              "SLA protocolos externos",
			Description:       "Meta para acompanhamento de protocolos e retornos a cidadãos.",
			FirstResponseTime: "4h",
			ResolutionTime:    "48h",
			NextResponseTime:  "12h",
			Notifications: []map[string]any{
				{"type": "warning", "time_delay_type": "before", "time_delay": "1h", "metric": "next_response", "recipients": []string{"assigned_user"}},
			},
		},
	} {
		id, err := ensureSLA(tx, item)
		if err != nil {
			return err
		}
		slaIDs[item.Name] = id
	}

	teamIDs := make(map[string]int64)
	for _, team := range []struct {
		Name              string
		Emoji             string
		BusinessHoursName string
		SLAName           string
	}{
		{Name: "Atendimento", Emoji: "🎧", BusinessHoursName: "Expediente padrão", SLAName: "SLA atendimento inicial"},
		{Name: "Ouvidoria", Emoji: "📝", BusinessHoursName: "Plantão da ouvidoria", SLAName: "SLA ouvidoria prioritária"},
		{Name: "Protocolos", Emoji: "📎", BusinessHoursName: "Cobertura integral", SLAName: "SLA protocolos externos"},
	} {
		id, err := ensureTeam(tx, team.Name, team.Emoji, businessHoursIDs[team.BusinessHoursName], slaIDs[team.SLAName])
		if err != nil {
			return err
		}
		teamIDs[team.Name] = id
	}

	inboxIDs := make(map[string]int64)
	for _, inbox := range []seedInbox{
		{Name: "Atendimento Gov", From: "atendimento@canalgov.local"},
		{Name: "Ouvidoria Gov", From: "ouvidoria@canalgov.local"},
		{Name: "Protocolos Gov", From: "protocolos@canalgov.local"},
	} {
		id, err := ensureInbox(tx, inbox.Name, inbox.From)
		if err != nil {
			return err
		}
		inboxIDs[inbox.Name] = id
	}

	agents := []seedAgent{
		{
			Email:     "ana.silva@canalgov.local",
			FirstName: "Ana",
			LastName:  "Silva",
			Roles:     []string{"Admin", "Gestor de Atendimento"},
			TeamName:  "Atendimento",
		},
		{
			Email:     "bruno.costa@canalgov.local",
			FirstName: "Bruno",
			LastName:  "Costa",
			Roles:     []string{"Agent", "Analista de Ouvidoria"},
			TeamName:  "Ouvidoria",
		},
		{
			Email:     "carla.souza@canalgov.local",
			FirstName: "Carla",
			LastName:  "Souza",
			Roles:     []string{"Agent", "Operador de Protocolo"},
			TeamName:  "Protocolos",
		},
	}

	agentIDs := make(map[string]int64)
	for _, agent := range agents {
		id, err := ensureAgent(tx, agent, password)
		if err != nil {
			return err
		}
		agentIDs[agent.Email] = id
		if err := ensureTeamMember(tx, teamIDs[agent.TeamName], id, "•"); err != nil {
			return err
		}
	}

	for _, tag := range []string{"urgente", "portal-gov", "beneficio", "cadastro", "sla", "intimacao", "vip", "documentacao"} {
		if _, err := ensureTag(tx, tag); err != nil {
			return err
		}
	}

	for _, attr := range []seedCustomAttribute{
		{Name: "Cartao SUS", Description: "Número do cartão SUS informado pelo cidadão.", AppliesTo: "contact", Key: "cartao_sus", DataType: "text"},
		{Name: "NIS", Description: "Número de Identificação Social para programas governamentais.", AppliesTo: "contact", Key: "nis", DataType: "text"},
		{Name: "Perfil do cidadão", Description: "Faixa de relacionamento predominante com o serviço.", AppliesTo: "contact", Key: "perfil_cidadao", DataType: "list", Values: []string{"beneficiário", "servidor", "fornecedor"}},
		{Name: "Canal de origem", Description: "Canal predominante de entrada do chamado.", AppliesTo: "conversation", Key: "canal_origem", DataType: "list", Values: []string{"portal", "email", "oficio"}},
		{Name: "Número do protocolo externo", Description: "Protocolo recebido de sistema ou órgão externo.", AppliesTo: "conversation", Key: "protocolo_externo", DataType: "text"},
		{Name: "Exige resposta formal", Description: "Define se a conversa exige retorno formal ao cidadão.", AppliesTo: "conversation", Key: "exige_resposta_formal", DataType: "checkbox"},
	} {
		if _, err := ensureCustomAttribute(tx, attr); err != nil {
			return err
		}
	}

	templates := []seedTemplate{
		{
			Name:      "Resposta padrão inicial",
			Type:      "response",
			Subject:   "Atualização do seu atendimento",
			Body:      "<p>Olá {{ .Conversation.Contact.FirstName }},</p><p>Recebemos sua solicitação e já iniciamos a análise. Se precisar complementar com documentos ou informações, responda esta mensagem.</p><p>Atenciosamente,<br>Equipe CanalGov</p>",
			IsDefault: true,
		},
		{
			Name:      "Retorno da ouvidoria",
			Type:      "response",
			Subject:   "Retorno da ouvidoria",
			Body:      "<p>Olá {{ .Conversation.Contact.FirstName }},</p><p>Sua manifestação foi encaminhada para análise da ouvidoria. Assim que a apuração for concluída, retornaremos com os próximos passos.</p><p>Atenciosamente,<br>Ouvidoria CanalGov</p>",
			TeamName:  "Ouvidoria",
			IsDefault: true,
		},
		{
			Name:      "Atualização de protocolo",
			Type:      "response",
			Subject:   "Atualização de protocolo",
			Body:      "<p>Olá {{ .Conversation.Contact.FirstName }},</p><p>Seu protocolo segue em tratamento pela equipe responsável. Se houver pendências, entraremos em contato pelos canais cadastrados.</p><p>Atenciosamente,<br>Equipe de Protocolos</p>",
			TeamName:  "Protocolos",
			IsDefault: true,
		},
		{
			Name:    "Confirmação de recebimento de documentos",
			Type:    "response",
			Subject: "Documentos recebidos com sucesso",
			Body:    "<p>Olá {{ .Conversation.Contact.FirstName }},</p><p>Confirmamos o recebimento dos documentos enviados. A equipe vai analisar o material e retornará com a próxima atualização.</p><p>Atenciosamente,<br>Equipe CanalGov</p>",
		},
		{
			Name:     "Solicitação de resposta formal",
			Type:     "response",
			Subject:  "Retorno formal em andamento",
			Body:     "<p>Olá {{ .Conversation.Contact.FirstName }},</p><p>Seu atendimento exige resposta formal. O posicionamento está em elaboração e será encaminhado por este canal.</p><p>Atenciosamente,<br>CanalGov</p>",
			TeamName: "Ouvidoria",
		},
		{
			Name:     "Pendência de protocolo externo",
			Type:     "response",
			Subject:  "Aguardando retorno do órgão responsável",
			Body:     "<p>Olá {{ .Conversation.Contact.FirstName }},</p><p>Estamos aguardando atualização do órgão responsável pelo protocolo externo. Assim que houver retorno, atualizaremos o chamado.</p><p>Atenciosamente,<br>Equipe de Protocolos</p>",
			TeamName: "Protocolos",
		},
	}
	for _, tpl := range templates {
		var teamID *int64
		if tpl.TeamName != "" {
			id := teamIDs[tpl.TeamName]
			teamID = &id
		}
		if _, err := ensureTemplate(tx, tpl, teamID); err != nil {
			return err
		}
	}

	for _, macro := range []seedMacro{
		{
			Name:           "Solicitar documentação complementar",
			MessageContent: "<p>Olá, precisamos de documentação complementar para dar andamento ao atendimento. Assim que possível, responda com os anexos necessários.</p>",
			Visibility:     "all",
			VisibleWhen:    []string{"replying", "starting_conversation"},
			Actions: []map[string]any{
				{"type": "add_tags", "value": []string{"cadastro"}},
			},
		},
		{
			Name:           "Encaminhar para ouvidoria",
			MessageContent: "<p>Registro encaminhado para análise da equipe de ouvidoria. Retornaremos assim que houver posicionamento.</p>",
			Visibility:     "team",
			VisibleWhen:    []string{"replying"},
			TeamName:       "Ouvidoria",
			Actions: []map[string]any{
				{"type": "assign_team", "value": []string{fmt.Sprintf("%d", teamIDs["Ouvidoria"])}},
				{"type": "set_priority", "value": []string{"2"}},
			},
		},
		{
			Name:           "Cobrar atualização de protocolo",
			MessageContent: "<p>Estamos acompanhando a atualização do seu protocolo e retornaremos com a resposta oficial assim que houver andamento.</p>",
			Visibility:     "user",
			VisibleWhen:    []string{"replying", "adding_private_note"},
			AgentEmail:     "carla.souza@canalgov.local",
			Actions: []map[string]any{
				{"type": "add_tags", "value": []string{"sla"}},
				{"type": "send_private_note", "value": []string{"Cobrança interna registrada para acompanhamento do protocolo."}},
			},
		},
		{
			Name:           "Confirmar recebimento de anexo",
			MessageContent: "<p>Recebemos o documento encaminhado e ele já foi anexado ao seu atendimento.</p>",
			Visibility:     "all",
			VisibleWhen:    []string{"replying"},
			Actions: []map[string]any{
				{"type": "add_tags", "value": []string{"documentacao"}},
			},
		},
		{
			Name:           "Registrar atendimento VIP",
			MessageContent: "<p>Seu atendimento foi priorizado e segue em acompanhamento direto da equipe responsável.</p>",
			Visibility:     "team",
			VisibleWhen:    []string{"replying", "adding_private_note"},
			TeamName:       "Atendimento",
			Actions: []map[string]any{
				{"type": "add_tags", "value": []string{"vip"}},
				{"type": "set_priority", "value": []string{"1"}},
			},
		},
		{
			Name:           "Intimação recebida",
			MessageContent: "<p>Recebemos a comunicação oficial e ela foi encaminhada para análise da equipe competente.</p>",
			Visibility:     "user",
			VisibleWhen:    []string{"replying"},
			AgentEmail:     "bruno.costa@canalgov.local",
			Actions: []map[string]any{
				{"type": "add_tags", "value": []string{"intimacao"}},
				{"type": "send_private_note", "value": []string{"Comunicação oficial recebida e classificada para resposta formal."}},
			},
		},
	} {
		if _, err := ensureMacro(tx, macro, teamIDs, agentIDs); err != nil {
			return err
		}
	}

	for _, webhook := range []seedWebhook{
		{Name: "Webhook atendimento", URL: "https://example.test/hooks/atendimento", Events: []string{"conversation.created", "message.created"}, Secret: "segredo-atendimento", IsActive: true},
		{Name: "Webhook ouvidoria", URL: "https://example.test/hooks/ouvidoria", Events: []string{"conversation.status_changed", "conversation.tags_changed"}, Secret: "segredo-ouvidoria", IsActive: true},
		{Name: "Webhook protocolos", URL: "https://example.test/hooks/protocolos", Events: []string{"conversation.assigned", "conversation.unassigned"}, Secret: "segredo-protocolos", IsActive: false},
		{Name: "Webhook mensagens", URL: "https://example.test/hooks/mensagens", Events: []string{"message.updated"}, Secret: "segredo-mensagens", IsActive: true},
		{Name: "Webhook sla", URL: "https://example.test/hooks/sla", Events: []string{"conversation.status_changed", "conversation.assigned"}, Secret: "segredo-sla", IsActive: true},
		{Name: "Webhook auditoria", URL: "https://example.test/hooks/auditoria", Events: []string{"conversation.created", "conversation.unassigned"}, Secret: "segredo-auditoria", IsActive: false},
	} {
		if _, err := ensureWebhook(tx, webhook); err != nil {
			return err
		}
	}

	contacts := []seedContact{
		{Key: "maria", Email: "maria.santos@example.com", FirstName: "Maria", LastName: "Santos", Phone: "11999990001", CustomAttributes: map[string]any{"cartao_sus": "12345622", "nis": "99887766554", "perfil_cidadao": "beneficiário"}},
		{Key: "joao", Email: "joao.oliveira@example.com", FirstName: "Joao", LastName: "Oliveira", Phone: "11999990002", CustomAttributes: map[string]any{"cartao_sus": "22334455", "perfil_cidadao": "servidor"}},
		{Key: "fernanda", Email: "fernanda.lima@example.com", FirstName: "Fernanda", LastName: "Lima", Phone: "11999990003", CustomAttributes: map[string]any{"nis": "11223344556", "perfil_cidadao": "fornecedor"}},
		{Key: "paulo", Email: "paulo.almeida@example.com", FirstName: "Paulo", LastName: "Almeida", Phone: "11999990004", CustomAttributes: map[string]any{"cartao_sus": "44556677"}},
		{Key: "lucia", Email: "lucia.ferreira@example.com", FirstName: "Lucia", LastName: "Ferreira", Phone: "11999990005", CustomAttributes: map[string]any{"nis": "55443322110"}},
		{Key: "roberto", Email: "roberto.gomes@example.com", FirstName: "Roberto", LastName: "Gomes", Phone: "11999990006", CustomAttributes: map[string]any{"cartao_sus": "77889900", "perfil_cidadao": "servidor"}},
		{Key: "aline", Email: "aline.barbosa@example.com", FirstName: "Aline", LastName: "Barbosa", Phone: "11999990007", CustomAttributes: map[string]any{"nis": "66778899001", "perfil_cidadao": "beneficiário"}},
		{Key: "diego", Email: "diego.ramos@example.com", FirstName: "Diego", LastName: "Ramos", Phone: "11999990008", CustomAttributes: map[string]any{"cartao_sus": "88990011", "perfil_cidadao": "fornecedor"}},
		{Key: "patricia", Email: "patricia.melo@example.com", FirstName: "Patricia", LastName: "Melo", Phone: "11999990009", CustomAttributes: map[string]any{"nis": "99001122334", "perfil_cidadao": "servidor"}},
	}

	contactIDs := make(map[string]int64)
	channelIDs := make(map[string]int64)
	for _, contact := range contacts {
		contactIDs[contact.Key] = 0
		for inboxName, inboxID := range inboxIDs {
			contactID, channelID, err := ensureContact(tx, inboxID, contact)
			if err != nil {
				return err
			}
			contactIDs[contact.Key] = contactID
			channelIDs[fmt.Sprintf("%s|%s", contact.Key, inboxName)] = channelID
		}
	}

	for _, note := range []struct {
		ContactKey string
		AgentEmail string
		Content    string
	}{
		{ContactKey: "maria", AgentEmail: "ana.silva@canalgov.local", Content: "Contato acompanha vacinação pelo portal e costuma responder no mesmo dia."},
		{ContactKey: "joao", AgentEmail: "bruno.costa@canalgov.local", Content: "Caso sensível ligado a benefício social. Priorizar clareza no retorno."},
		{ContactKey: "fernanda", AgentEmail: "carla.souza@canalgov.local", Content: "Relata erro intermitente ao atualizar cadastro em horário de pico."},
		{ContactKey: "aline", AgentEmail: "ana.silva@canalgov.local", Content: "Solicita atualização frequente e prefere retorno por e-mail."},
		{ContactKey: "diego", AgentEmail: "carla.souza@canalgov.local", Content: "Fornecedor com pendência documental recorrente."},
		{ContactKey: "patricia", AgentEmail: "bruno.costa@canalgov.local", Content: "Atendimento com potencial de resposta formal e acompanhamento jurídico."},
	} {
		if err := ensureContactNote(tx, contactIDs[note.ContactKey], agentIDs[note.AgentEmail], note.Content); err != nil {
			return err
		}
	}

	statusIDs, err := loadLookup(tx, "conversation_statuses")
	if err != nil {
		return err
	}
	priorityIDs, err := loadLookup(tx, "conversation_priorities")
	if err != nil {
		return err
	}

	conversations := []seedConversation{
		{
			Key:              "vacina-covid",
			Subject:          "Duvida sobre agendamento de vacina",
			ContactKey:       "maria",
			InboxName:        "Atendimento Gov",
			Status:           "Open",
			Priority:         "High",
			TeamName:         "Atendimento",
			AssigneeEmail:    "ana.silva@canalgov.local",
			TagNames:         []string{"urgente", "portal-gov"},
			CreatedAgo:       2 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "portal", "protocolo_externo": "VAC-2026-0001", "exige_resposta_formal": false},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "maria.santos@example.com", Content: "Bom dia, nao estou conseguindo concluir o agendamento da vacina pelo portal.", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "ana.silva@canalgov.local", Content: "Bom dia, Maria. Vou verificar com voce. Pode me informar o CPF e o navegador utilizado?", Offset: 12 * time.Minute},
			},
		},
		{
			Key:              "auxilio-moradia",
			Subject:          "Solicitacao de revisao do beneficio",
			ContactKey:       "joao",
			InboxName:        "Ouvidoria Gov",
			Status:           "Resolved",
			Priority:         "Medium",
			TeamName:         "Ouvidoria",
			AssigneeEmail:    "bruno.costa@canalgov.local",
			TagNames:         []string{"beneficio", "sla"},
			CreatedAgo:       8 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "email", "protocolo_externo": "OUV-2026-0042", "exige_resposta_formal": true},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "joao.oliveira@example.com", Content: "Gostaria de revisar o valor do meu beneficio, pois houve alteracao na renda familiar.", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "bruno.costa@canalgov.local", Content: "Recebemos sua solicitacao e registramos para analise da equipe responsavel.", Offset: 25 * time.Minute},
				{Key: "m3", Type: "outgoing", SenderType: "agent", SenderEmail: "bruno.costa@canalgov.local", Content: "Analise concluida. O beneficio foi recalculado e a atualizacao sera refletida no proximo ciclo.", Offset: 90 * time.Minute},
			},
		},
		{
			Key:              "erro-cadastro",
			Subject:          "Erro ao atualizar cadastro no portal",
			ContactKey:       "fernanda",
			InboxName:        "Protocolos Gov",
			Status:           "Open",
			Priority:         "Medium",
			TeamName:         "Protocolos",
			AssigneeEmail:    "carla.souza@canalgov.local",
			TagNames:         []string{"cadastro", "portal-gov"},
			CreatedAgo:       26 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "portal", "protocolo_externo": "CAD-2026-0188", "exige_resposta_formal": false},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "fernanda.lima@example.com", Content: "Ao salvar meu cadastro, o sistema retorna erro interno.", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "carla.souza@canalgov.local", Content: "Estamos analisando o incidente. Se possivel, envie uma captura de tela com a mensagem exibida.", Offset: 18 * time.Minute},
			},
		},
		{
			Key:              "protocolo-atrasado",
			Subject:          "Acompanhamento de protocolo sem retorno",
			ContactKey:       "paulo",
			InboxName:        "Protocolos Gov",
			Status:           "Snoozed",
			Priority:         "Low",
			TeamName:         "Protocolos",
			AssigneeEmail:    "carla.souza@canalgov.local",
			TagNames:         []string{"sla"},
			CreatedAgo:       48 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "oficio", "protocolo_externo": "PRO-2026-3001", "exige_resposta_formal": true},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "paulo.almeida@example.com", Content: "Meu protocolo segue parado desde a semana passada. Preciso de uma previsao.", Offset: 0},
			},
		},
		{
			Key:              "segunda-via",
			Subject:          "Pedido de segunda via de documento",
			ContactKey:       "lucia",
			InboxName:        "Atendimento Gov",
			Status:           "Closed",
			Priority:         "Low",
			TeamName:         "Atendimento",
			AssigneeEmail:    "ana.silva@canalgov.local",
			TagNames:         []string{"cadastro"},
			CreatedAgo:       72 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "email", "protocolo_externo": "DOC-2026-1100", "exige_resposta_formal": false},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "lucia.ferreira@example.com", Content: "Preciso emitir a segunda via do comprovante e nao achei a opcao no portal.", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "ana.silva@canalgov.local", Content: "A emissao pode ser feita pelo menu Servicos > Documentos. Encaminhei tambem o passo a passo por email.", Offset: 14 * time.Minute},
			},
		},
		{
			Key:              "oficio-judicial",
			Subject:          "Recebimento de ofício judicial para manifestação",
			ContactKey:       "patricia",
			InboxName:        "Ouvidoria Gov",
			Status:           "Open",
			Priority:         "High",
			TeamName:         "Ouvidoria",
			AssigneeEmail:    "bruno.costa@canalgov.local",
			TagNames:         []string{"intimacao", "vip"},
			CreatedAgo:       90 * time.Minute,
			CustomAttributes: map[string]any{"canal_origem": "oficio", "protocolo_externo": "JUD-2026-7781", "exige_resposta_formal": true},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "patricia.melo@example.com", Content: "Encaminho ofício judicial com prazo de manifestação em 48 horas.", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "bruno.costa@canalgov.local", Content: "Ofício recebido. A equipe responsável já iniciou a análise para elaboração da resposta formal.", Offset: 20 * time.Minute},
			},
		},
		{
			Key:              "pendencia-documental",
			Subject:          "Pendência de documentação para cadastro de fornecedor",
			ContactKey:       "diego",
			InboxName:        "Protocolos Gov",
			Status:           "Open",
			Priority:         "Medium",
			TeamName:         "Protocolos",
			AssigneeEmail:    "carla.souza@canalgov.local",
			TagNames:         []string{"documentacao", "cadastro"},
			CreatedAgo:       6 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "email", "protocolo_externo": "FOR-2026-2230", "exige_resposta_formal": false},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "diego.ramos@example.com", Content: "Recebi exigência de documentação complementar para cadastro de fornecedor. Quais arquivos faltam?", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "carla.souza@canalgov.local", Content: "Faltam a certidão negativa atualizada e o comprovante de endereço empresarial.", Offset: 40 * time.Minute},
			},
		},
		{
			Key:              "atendimento-vip",
			Subject:          "Solicitação prioritária de acesso ao portal interno",
			ContactKey:       "roberto",
			InboxName:        "Atendimento Gov",
			Status:           "Replied",
			Priority:         "High",
			TeamName:         "Atendimento",
			AssigneeEmail:    "ana.silva@canalgov.local",
			TagNames:         []string{"vip", "portal-gov"},
			CreatedAgo:       10 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "portal", "protocolo_externo": "VIP-2026-0099", "exige_resposta_formal": false},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "roberto.gomes@example.com", Content: "Preciso de acesso urgente ao portal interno para fechamento do mês.", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "ana.silva@canalgov.local", Content: "Seu acesso foi regularizado. Orientei também a atualização da senha e da autenticação em duas etapas.", Offset: 25 * time.Minute},
			},
		},
		{
			Key:              "beneficio-bloqueado",
			Subject:          "Benefício bloqueado após atualização cadastral",
			ContactKey:       "aline",
			InboxName:        "Ouvidoria Gov",
			Status:           "Open",
			Priority:         "High",
			TeamName:         "Ouvidoria",
			AssigneeEmail:    "bruno.costa@canalgov.local",
			TagNames:         []string{"beneficio", "urgente"},
			CreatedAgo:       14 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "portal", "protocolo_externo": "BEN-2026-3120", "exige_resposta_formal": true},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "aline.barbosa@example.com", Content: "Meu benefício foi bloqueado depois da atualização cadastral e preciso de análise urgente.", Offset: 0},
				{Key: "m2", Type: "outgoing", SenderType: "agent", SenderEmail: "bruno.costa@canalgov.local", Content: "Abrimos análise prioritária do bloqueio e pedimos verificação documental à área responsável.", Offset: 35 * time.Minute},
			},
		},
		{
			Key:              "certidao-servidor",
			Subject:          "Emissão de certidão funcional com erro",
			ContactKey:       "patricia",
			InboxName:        "Atendimento Gov",
			Status:           "Snoozed",
			Priority:         "Medium",
			TeamName:         "Atendimento",
			AssigneeEmail:    "ana.silva@canalgov.local",
			TagNames:         []string{"cadastro", "documentacao"},
			CreatedAgo:       30 * time.Hour,
			CustomAttributes: map[string]any{"canal_origem": "portal", "protocolo_externo": "CER-2026-4110", "exige_resposta_formal": false},
			Messages: []seedMessage{
				{Key: "m1", Type: "incoming", SenderType: "contact", SenderEmail: "patricia.melo@example.com", Content: "Ao tentar emitir a certidão funcional, o sistema informa inconsistência cadastral.", Offset: 0},
			},
		},
	}

	for _, conv := range conversations {
		inboxID := inboxIDs[conv.InboxName]
		channelKey := fmt.Sprintf("%s|%s", conv.ContactKey, conv.InboxName)
		convID, err := ensureConversation(tx, conv, contactIDs[conv.ContactKey], channelIDs[channelKey], inboxID, teamIDs[conv.TeamName], agentIDs[conv.AssigneeEmail], statusIDs[conv.Status], priorityIDs[conv.Priority])
		if err != nil {
			return err
		}

		for _, tagName := range conv.TagNames {
			tagID, err := ensureTag(tx, tagName)
			if err != nil {
				return err
			}
			if err := ensureConversationTag(tx, convID, tagID); err != nil {
				return err
			}
		}

		if err := ensureMessages(tx, conv, convID, contactIDs, agentIDs); err != nil {
			return err
		}

		if err := refreshConversationState(tx, conv, convID); err != nil {
			return err
		}
	}

	for _, view := range []seedView{
		{Name: "Atendimento urgente", Visibility: "all", Filters: []map[string]any{{"model": "conversations", "field": "priority_id", "operator": "equals", "value": fmt.Sprintf("%d", priorityIDs["High"])}}},
		{Name: "Fila de ouvidoria", Visibility: "team", TeamName: "Ouvidoria", Filters: []map[string]any{{"model": "conversations", "field": "assigned_team_id", "operator": "equals", "value": fmt.Sprintf("%d", teamIDs["Ouvidoria"])}}},
		{Name: "Minhas conversas abertas", Visibility: "user", AgentEmail: "ana.silva@canalgov.local", Filters: []map[string]any{{"model": "conversations", "field": "assigned_user_id", "operator": "equals", "value": fmt.Sprintf("%d", agentIDs["ana.silva@canalgov.local"])}, {"model": "conversations", "field": "status_id", "operator": "not equals", "value": fmt.Sprintf("%d", statusIDs["Closed"])}}},
		{Name: "Protocolos críticos", Visibility: "team", TeamName: "Protocolos", Filters: []map[string]any{{"model": "conversations", "field": "priority_id", "operator": "equals", "value": fmt.Sprintf("%d", priorityIDs["High"])}, {"model": "conversations", "field": "assigned_team_id", "operator": "equals", "value": fmt.Sprintf("%d", teamIDs["Protocolos"])}}},
		{Name: "Benefícios em acompanhamento", Visibility: "all", Filters: []map[string]any{{"model": "conversations", "field": "status_id", "operator": "not equals", "value": fmt.Sprintf("%d", statusIDs["Resolved"])}}},
		{Name: "Conversas do Bruno", Visibility: "user", AgentEmail: "bruno.costa@canalgov.local", Filters: []map[string]any{{"model": "conversations", "field": "assigned_user_id", "operator": "equals", "value": fmt.Sprintf("%d", agentIDs["bruno.costa@canalgov.local"])}}},
	} {
		if _, err := ensureView(tx, view, teamIDs, agentIDs); err != nil {
			return err
		}
	}

	for _, rule := range []seedAutomationRule{
		{
			Name:          "Nova conversa urgente",
			Description:   "Marca novas conversas com referência a protocolo como prioritárias.",
			Type:          "new_conversation",
			Enabled:       true,
			Weight:        10,
			ExecutionMode: "all",
			Rules: []map[string]any{
				{
					"group_operator": "OR",
					"groups": []map[string]any{
						{"logical_op": "OR", "rules": []map[string]any{{"field": "subject", "field_type": "conversation", "operator": "contains", "value": "protocolo", "case_sensitive_match": false}}},
						{"logical_op": "OR", "rules": []map[string]any{}},
					},
					"actions": []map[string]any{{"type": "set_priority", "value": []string{fmt.Sprintf("%d", priorityIDs["High"])}}},
				},
			},
		},
		{
			Name:          "Mudança de status em ouvidoria",
			Description:   "Quando a conversa da ouvidoria for respondida, aplica tag de acompanhamento.",
			Type:          "conversation_update",
			Events:        []string{"conversation.status.change"},
			Enabled:       true,
			Weight:        20,
			ExecutionMode: "all",
			Rules: []map[string]any{
				{
					"group_operator": "AND",
					"groups": []map[string]any{
						{"logical_op": "OR", "rules": []map[string]any{{"field": "assigned_team", "field_type": "conversation", "operator": "equals", "value": fmt.Sprintf("%d", teamIDs["Ouvidoria"]), "case_sensitive_match": false}}},
						{"logical_op": "OR", "rules": []map[string]any{{"field": "status", "field_type": "conversation", "operator": "equals", "value": "Replied", "case_sensitive_match": false}}},
					},
					"actions": []map[string]any{{"type": "add_tags", "value": []string{"beneficio"}}},
				},
			},
		},
		{
			Name:          "Tempo de espera elevado",
			Description:   "Aplica nota privada para conversas com muitas horas sem resposta.",
			Type:          "time_trigger",
			Enabled:       true,
			Weight:        30,
			ExecutionMode: "first_match",
			Rules: []map[string]any{
				{
					"group_operator": "OR",
					"groups": []map[string]any{
						{"logical_op": "OR", "rules": []map[string]any{{"field": "hours_since_created", "field_type": "conversation", "operator": "greater than", "value": "24", "case_sensitive_match": false}}},
						{"logical_op": "OR", "rules": []map[string]any{}},
					},
					"actions": []map[string]any{{"type": "send_private_note", "value": []string{"Conversa em acompanhamento automático por tempo de espera elevado."}}},
				},
			},
		},
	} {
		if _, err := ensureAutomationRule(tx, rule); err != nil {
			return err
		}
	}

	return nil
}

func ensureInbox(tx *sqlx.Tx, name, from string) (int64, error) {
	var id int64
	err := tx.Get(&id, `SELECT id FROM inboxes WHERE name = $1 LIMIT 1`, name)
	if err == nil {
		_, err = tx.Exec(`UPDATE inboxes SET "from" = $2, enabled = true, updated_at = now() WHERE id = $1`, id, from)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	err = tx.Get(&id, `
		INSERT INTO inboxes (name, channel, enabled, "from", config)
		VALUES ($1, 'email', true, $2, '{}'::jsonb)
		RETURNING id
	`, name, from)
	return id, err
}

func ensureBusinessHours(tx *sqlx.Tx, item seedBusinessHours) (int64, error) {
	var (
		id        int64
		hoursJSON []byte
	)
	hoursJSON, _ = json.Marshal(item.Hours)
	err := tx.Get(&id, `SELECT id FROM business_hours WHERE name = $1`, item.Name)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE business_hours
			SET description = $2, is_always_open = $3, hours = $4::jsonb, holidays = '{}'::jsonb, updated_at = now()
			WHERE id = $1
		`, id, item.Description, item.IsAlwaysOpen, string(hoursJSON))
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `
		INSERT INTO business_hours (name, description, is_always_open, hours, holidays)
		VALUES ($1, $2, $3, $4::jsonb, '{}'::jsonb)
		RETURNING id
	`, item.Name, item.Description, item.IsAlwaysOpen, string(hoursJSON))
	return id, err
}

func ensureSLA(tx *sqlx.Tx, item seedSLA) (int64, error) {
	var (
		id                int64
		notificationsJSON []byte
	)
	notificationsJSON, _ = json.Marshal(item.Notifications)
	err := tx.Get(&id, `SELECT id FROM sla_policies WHERE name = $1`, item.Name)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE sla_policies
			SET description = $2,
			    first_response_time = $3,
			    resolution_time = $4,
			    next_response_time = $5,
			    notifications = $6::jsonb,
			    updated_at = now()
			WHERE id = $1
		`, id, item.Description, item.FirstResponseTime, item.ResolutionTime, item.NextResponseTime, string(notificationsJSON))
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `
		INSERT INTO sla_policies (name, description, first_response_time, resolution_time, next_response_time, notifications)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb)
		RETURNING id
	`, item.Name, item.Description, item.FirstResponseTime, item.ResolutionTime, item.NextResponseTime, string(notificationsJSON))
	return id, err
}

func ensureTeam(tx *sqlx.Tx, name, emoji string, businessHoursID, slaID int64) (int64, error) {
	var id int64
	err := tx.Get(&id, `SELECT id FROM teams WHERE name = $1`, name)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE teams
			SET emoji = $2,
			    conversation_assignment_type = 'Manual',
			    timezone = 'America/Sao_Paulo',
			    business_hours_id = $3,
			    sla_policy_id = $4,
			    updated_at = now()
			WHERE id = $1
		`, id, emoji, businessHoursID, slaID)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	err = tx.Get(&id, `
		INSERT INTO teams (name, emoji, conversation_assignment_type, timezone, business_hours_id, sla_policy_id)
		VALUES ($1, $2, 'Manual', 'America/Sao_Paulo', $3, $4)
		RETURNING id
	`, name, emoji, businessHoursID, slaID)
	return id, err
}

func ensureAgent(tx *sqlx.Tx, agent seedAgent, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int64
	err = tx.Get(&id, `SELECT id FROM users WHERE email = $1 AND type = 'agent' AND deleted_at IS NULL`, agent.Email)
	if err == nil {
		if _, err := tx.Exec(`
			UPDATE users
			SET first_name = $2, last_name = $3, password = $4, enabled = true, updated_at = now()
			WHERE id = $1
		`, id, agent.FirstName, agent.LastName, hash); err != nil {
			return 0, err
		}
	} else if err == sql.ErrNoRows {
		if err := tx.Get(&id, `
			INSERT INTO users (email, type, first_name, last_name, password, enabled, availability_status)
			VALUES ($1, 'agent', $2, $3, $4, true, 'online')
			RETURNING id
		`, agent.Email, agent.FirstName, agent.LastName, hash); err != nil {
			return 0, err
		}
	} else {
		return 0, err
	}

	if _, err := tx.Exec(`DELETE FROM user_roles WHERE user_id = $1`, id); err != nil {
		return 0, err
	}

	for _, roleName := range agent.Roles {
		if _, err := tx.Exec(`
			INSERT INTO user_roles (user_id, role_id)
			SELECT $1, id FROM roles WHERE name = $2
			ON CONFLICT (user_id, role_id) DO NOTHING
		`, id, roleName); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func ensureRole(tx *sqlx.Tx, role seedRole) (int64, error) {
	var id int64
	err := tx.Get(&id, `SELECT id FROM roles WHERE name = $1`, role.Name)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE roles
			SET description = $2, permissions = $3, updated_at = now()
			WHERE id = $1
		`, id, role.Description, pq.StringArray(role.Permissions))
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	err = tx.Get(&id, `
		INSERT INTO roles (name, description, permissions)
		VALUES ($1, $2, $3)
		RETURNING id
	`, role.Name, role.Description, pq.StringArray(role.Permissions))
	return id, err
}

func ensureTemplate(tx *sqlx.Tx, tpl seedTemplate, teamID *int64) (int64, error) {
	var (
		id          int64
		teamIDValue any
	)
	if teamID != nil {
		teamIDValue = *teamID
	}

	err := tx.Get(&id, `
		SELECT id
		FROM templates
		WHERE name = $1 AND type = $2 AND COALESCE(team_id, 0) = COALESCE($3, 0)
		LIMIT 1
	`, tpl.Name, tpl.Type, teamIDValue)
	if err == nil {
		if _, err := tx.Exec(`
			UPDATE templates
			SET subject = $2,
			    body = $3,
			    is_default = $4,
			    is_builtin = false,
			    team_id = $5,
			    updated_at = now()
			WHERE id = $1
		`, id, tpl.Subject, tpl.Body, tpl.IsDefault, teamIDValue); err != nil {
			return 0, err
		}
	} else if err == sql.ErrNoRows {
		if err := tx.Get(&id, `
			INSERT INTO templates (type, name, subject, body, is_default, is_builtin, team_id)
			VALUES ($1, $2, $3, $4, $5, false, $6)
			RETURNING id
		`, tpl.Type, tpl.Name, tpl.Subject, tpl.Body, tpl.IsDefault, teamIDValue); err != nil {
			return 0, err
		}
	} else {
		return 0, err
	}

	if tpl.IsDefault {
		if _, err := tx.Exec(`
			UPDATE templates
			SET is_default = false, updated_at = now()
			WHERE type = $1
			  AND COALESCE(team_id, 0) = COALESCE($2, 0)
			  AND id <> $3
		`, tpl.Type, teamIDValue, id); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func ensureTeamMember(tx *sqlx.Tx, teamID, userID int64, emoji string) error {
	_, err := tx.Exec(`
		INSERT INTO team_members (team_id, user_id, emoji)
		VALUES ($1, $2, $3)
		ON CONFLICT (team_id, user_id) DO UPDATE SET emoji = EXCLUDED.emoji, updated_at = now()
	`, teamID, userID, emoji)
	return err
}

func ensureTag(tx *sqlx.Tx, name string) (int64, error) {
	var id int64
	err := tx.Get(&id, `SELECT id FROM tags WHERE name = $1`, name)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `INSERT INTO tags (name) VALUES ($1) RETURNING id`, name)
	return id, err
}

func ensureCustomAttribute(tx *sqlx.Tx, attr seedCustomAttribute) (int64, error) {
	var id int64
	values := pq.StringArray(attr.Values)
	if values == nil {
		values = pq.StringArray{}
	}
	err := tx.Get(&id, `SELECT id FROM custom_attribute_definitions WHERE key = $1 AND applies_to = $2`, attr.Key, attr.AppliesTo)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE custom_attribute_definitions
			SET name = $2,
			    description = $3,
			    values = $4,
			    data_type = $5,
			    regex = $6,
			    regex_hint = $7,
			    updated_at = now()
			WHERE id = $1
		`, id, attr.Name, attr.Description, values, attr.DataType, attr.Regex, attr.RegexHint)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `
		INSERT INTO custom_attribute_definitions (applies_to, name, description, key, values, data_type, regex, regex_hint)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, attr.AppliesTo, attr.Name, attr.Description, attr.Key, values, attr.DataType, attr.Regex, attr.RegexHint)
	return id, err
}

func ensureMacro(tx *sqlx.Tx, macro seedMacro, teamIDs, agentIDs map[string]int64) (int64, error) {
	var (
		id          int64
		actionsJSON []byte
		teamIDValue any
		userIDValue any
	)
	if macro.TeamName != "" {
		teamIDValue = teamIDs[macro.TeamName]
	}
	if macro.AgentEmail != "" {
		userIDValue = agentIDs[macro.AgentEmail]
	}
	actionsJSON, _ = json.Marshal(macro.Actions)
	err := tx.Get(&id, `SELECT id FROM macros WHERE name = $1 LIMIT 1`, macro.Name)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE macros
			SET message_content = $2,
			    visibility = $3,
			    visible_when = $4,
			    actions = $5::jsonb,
			    team_id = $6,
			    user_id = $7,
			    updated_at = now()
			WHERE id = $1
		`, id, macro.MessageContent, macro.Visibility, pq.StringArray(macro.VisibleWhen), string(actionsJSON), teamIDValue, userIDValue)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `
		INSERT INTO macros (name, actions, visibility, visible_when, message_content, user_id, team_id)
		VALUES ($1, $2::jsonb, $3, $4, $5, $6, $7)
		RETURNING id
	`, macro.Name, string(actionsJSON), macro.Visibility, pq.StringArray(macro.VisibleWhen), macro.MessageContent, userIDValue, teamIDValue)
	return id, err
}

func ensureWebhook(tx *sqlx.Tx, webhook seedWebhook) (int64, error) {
	var id int64
	err := tx.Get(&id, `SELECT id FROM webhooks WHERE name = $1 LIMIT 1`, webhook.Name)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE webhooks
			SET url = $2, events = $3, secret = $4, is_active = $5, updated_at = now()
			WHERE id = $1
		`, id, webhook.URL, pq.StringArray(webhook.Events), webhook.Secret, webhook.IsActive)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `
		INSERT INTO webhooks (name, url, events, secret, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, webhook.Name, webhook.URL, pq.StringArray(webhook.Events), webhook.Secret, webhook.IsActive)
	return id, err
}

func ensureView(tx *sqlx.Tx, view seedView, teamIDs, agentIDs map[string]int64) (int64, error) {
	var (
		id          int64
		filtersJSON []byte
		userIDValue any
		teamIDValue any
	)
	if view.AgentEmail != "" {
		userIDValue = agentIDs[view.AgentEmail]
	}
	if view.TeamName != "" {
		teamIDValue = teamIDs[view.TeamName]
	}
	filtersJSON, _ = json.Marshal(view.Filters)
	err := tx.Get(&id, `
		SELECT id FROM views
		WHERE name = $1
		  AND visibility = $2
		  AND COALESCE(user_id, 0) = COALESCE($3, 0)
		  AND COALESCE(team_id, 0) = COALESCE($4, 0)
		LIMIT 1
	`, view.Name, view.Visibility, userIDValue, teamIDValue)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE views
			SET filters = $2::jsonb, updated_at = now()
			WHERE id = $1
		`, id, string(filtersJSON))
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `
		INSERT INTO views (name, filters, visibility, user_id, team_id)
		VALUES ($1, $2::jsonb, $3, $4, $5)
		RETURNING id
	`, view.Name, string(filtersJSON), view.Visibility, userIDValue, teamIDValue)
	return id, err
}

func ensureAutomationRule(tx *sqlx.Tx, rule seedAutomationRule) (int64, error) {
	var (
		id        int64
		rulesJSON []byte
	)
	events := pq.StringArray(rule.Events)
	if events == nil {
		events = pq.StringArray{}
	}
	rulesJSON, _ = json.Marshal(rule.Rules)
	err := tx.Get(&id, `SELECT id FROM automation_rules WHERE name = $1 LIMIT 1`, rule.Name)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE automation_rules
			SET description = $2,
			    type = $3,
			    events = $4,
			    rules = $5::jsonb,
			    enabled = $6,
			    weight = $7,
			    execution_mode = $8,
			    updated_at = now()
			WHERE id = $1
		`, id, rule.Description, rule.Type, events, string(rulesJSON), rule.Enabled, rule.Weight, rule.ExecutionMode)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	err = tx.Get(&id, `
		INSERT INTO automation_rules (name, description, type, events, rules, enabled, weight, execution_mode)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7, $8)
		RETURNING id
	`, rule.Name, rule.Description, rule.Type, events, string(rulesJSON), rule.Enabled, rule.Weight, rule.ExecutionMode)
	return id, err
}

func ensureContact(tx *sqlx.Tx, inboxID int64, contact seedContact) (int64, int64, error) {
	var contactID int64
	customAttrs, _ := json.Marshal(contact.CustomAttributes)
	err := tx.Get(&contactID, `SELECT id FROM users WHERE email = $1 AND type = 'contact' AND deleted_at IS NULL`, contact.Email)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE users
			SET first_name = $2, last_name = $3, phone_number = $4, phone_number_country_code = '55', custom_attributes = $5::jsonb, updated_at = now()
			WHERE id = $1
		`, contactID, contact.FirstName, contact.LastName, contact.Phone, string(customAttrs))
		if err != nil {
			return 0, 0, err
		}
	} else if err == sql.ErrNoRows {
		err = tx.Get(&contactID, `
			INSERT INTO users (email, type, first_name, last_name, phone_number, phone_number_country_code, custom_attributes)
			VALUES ($1, 'contact', $2, $3, $4, '55', $5::jsonb)
			RETURNING id
		`, contact.Email, contact.FirstName, contact.LastName, contact.Phone, string(customAttrs))
		if err != nil {
			return 0, 0, err
		}
	} else {
		return 0, 0, err
	}

	identifier := contact.Email
	var channelID int64
	err = tx.Get(&channelID, `SELECT id FROM contact_channels WHERE contact_id = $1 AND inbox_id = $2`, contactID, inboxID)
	if err == nil {
		_, err = tx.Exec(`UPDATE contact_channels SET identifier = $2, updated_at = now() WHERE id = $1`, channelID, identifier)
		return contactID, channelID, err
	}
	if err != sql.ErrNoRows {
		return 0, 0, err
	}

	err = tx.Get(&channelID, `
		INSERT INTO contact_channels (contact_id, inbox_id, identifier)
		VALUES ($1, $2, $3)
		RETURNING id
	`, contactID, inboxID, identifier)
	return contactID, channelID, err
}

func ensureContactNote(tx *sqlx.Tx, contactID, userID int64, note string) error {
	var id int64
	err := tx.Get(&id, `SELECT id FROM contact_notes WHERE contact_id = $1 AND user_id = $2 AND note = $3 LIMIT 1`, contactID, userID, note)
	if err == nil {
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}
	_, err = tx.Exec(`INSERT INTO contact_notes (contact_id, note, user_id) VALUES ($1, $2, $3)`, contactID, note, userID)
	return err
}

func loadLookup(tx *sqlx.Tx, table string) (map[string]int64, error) {
	rows, err := tx.Queryx(fmt.Sprintf(`SELECT id, name FROM %s`, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]int64)
	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		out[name] = id
	}
	return out, rows.Err()
}

func ensureConversation(tx *sqlx.Tx, conv seedConversation, contactID, channelID, inboxID, teamID, assigneeID, statusID, priorityID int64) (int64, error) {
	var id int64
	err := tx.Get(&id, `
		SELECT id
		FROM conversations
		WHERE subject = $1 AND contact_id = $2 AND inbox_id = $3
		LIMIT 1
	`, conv.Subject, contactID, inboxID)
	createdAt := time.Now().Add(-conv.CreatedAgo)
	meta, _ := json.Marshal(map[string]any{})
	customAttrs, _ := json.Marshal(conv.CustomAttributes)
	if err == nil {
		_, err = tx.Exec(`
			UPDATE conversations
			SET subject = $2,
			    contact_id = $3,
			    contact_channel_id = $4,
			    inbox_id = $5,
			    assigned_team_id = $6,
			    assigned_user_id = $7,
			    status_id = $8,
			    priority_id = $9,
			    meta = $10,
			    custom_attributes = $11,
			    updated_at = now()
			WHERE id = $1
		`, id, conv.Subject, contactID, channelID, inboxID, teamID, assigneeID, statusID, priorityID, string(meta), string(customAttrs))
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	err = tx.Get(&id, `
		INSERT INTO conversations (
			created_at, updated_at, contact_id, assigned_user_id, assigned_team_id,
			inbox_id, contact_channel_id, status_id, priority_id, subject, meta,
			custom_attributes, waiting_since
		)
		VALUES (
			$1, $1, $2, $3, $4,
			$5, $6, $7, $8, $9, $10::jsonb,
			$11::jsonb, $1
		)
		RETURNING id
	`, createdAt, contactID, assigneeID, teamID, inboxID, channelID, statusID, priorityID, conv.Subject, string(meta), string(customAttrs))
	return id, err
}

func ensureConversationTag(tx *sqlx.Tx, conversationID, tagID int64) error {
	_, err := tx.Exec(`
		INSERT INTO conversation_tags (conversation_id, tag_id)
		VALUES ($1, $2)
		ON CONFLICT (conversation_id, tag_id) DO NOTHING
	`, conversationID, tagID)
	return err
}

func ensureMessages(tx *sqlx.Tx, conv seedConversation, conversationID int64, contactIDs, agentIDs map[string]int64) error {
	base := time.Now().Add(-conv.CreatedAgo)
	for _, msg := range conv.Messages {
		sourceID := fmt.Sprintf("dev-seed:%s:%s", conv.Key, msg.Key)
		var exists int64
		err := tx.Get(&exists, `SELECT id FROM conversation_messages WHERE source_id = $1 LIMIT 1`, sourceID)
		if err == nil {
			continue
		}
		if err != sql.ErrNoRows {
			return err
		}

		var senderID int64
		switch msg.SenderType {
		case "contact":
			senderID = contactIDs[conv.ContactKey]
		case "agent":
			senderID = agentIDs[msg.SenderEmail]
		default:
			return fmt.Errorf("sender_type invalido: %s", msg.SenderType)
		}

		createdAt := base.Add(msg.Offset)
		if _, err := tx.Exec(`
			INSERT INTO conversation_messages (
				created_at, updated_at, type, status, private, conversation_id,
				content_type, content, text_content, source_id, sender_id, sender_type, meta
			)
			VALUES (
				$1, $1, $2, 'sent', false, $3,
				'text', $4, $4, $5, $6, $7, '{}'::jsonb
			)
		`, createdAt, msg.Type, conversationID, msg.Content, sourceID, senderID, msg.SenderType); err != nil {
			return err
		}
	}
	return nil
}

func refreshConversationState(tx *sqlx.Tx, conv seedConversation, conversationID int64) error {
	type msgState struct {
		Content    string    `db:"content"`
		SenderType string    `db:"sender_type"`
		CreatedAt  time.Time `db:"created_at"`
	}

	var last msgState
	if err := tx.Get(&last, `
		SELECT content, sender_type, created_at
		FROM conversation_messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, conversationID); err != nil {
		return err
	}

	var firstReplyAt sql.NullTime
	var lastReplyAt sql.NullTime
	if err := tx.Get(&firstReplyAt, `
		SELECT min(created_at)
		FROM conversation_messages
		WHERE conversation_id = $1 AND sender_type = 'agent'
	`, conversationID); err != nil {
		return err
	}
	if err := tx.Get(&lastReplyAt, `
		SELECT max(created_at)
		FROM conversation_messages
		WHERE conversation_id = $1 AND sender_type = 'agent'
	`, conversationID); err != nil {
		return err
	}

	var snoozedUntil sql.NullTime
	if conv.Status == "Snoozed" {
		snoozedUntil = sql.NullTime{Time: time.Now().Add(6 * time.Hour), Valid: true}
	}

	var resolvedAt sql.NullTime
	if conv.Status == "Resolved" || conv.Status == "Closed" {
		resolvedAt = sql.NullTime{Time: last.CreatedAt, Valid: true}
	}

	var closedAt sql.NullTime
	if conv.Status == "Closed" {
		closedAt = sql.NullTime{Time: last.CreatedAt, Valid: true}
	}

	_, err := tx.Exec(`
		UPDATE conversations
		SET first_reply_at = $2,
		    last_reply_at = $3,
		    resolved_at = $4,
		    closed_at = $5,
		    snoozed_until = $6,
		    last_message_at = $7,
		    last_message = $8,
		    last_message_sender = $9,
		    last_interaction_at = $7,
		    last_interaction = $8,
		    last_interaction_sender = $9,
		    updated_at = now()
		WHERE id = $1
	`, conversationID, firstReplyAt, lastReplyAt, resolvedAt, closedAt, snoozedUntil, last.CreatedAt, last.Content, last.SenderType)
	return err
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("dev-seed: ")
}
