-- Add new channels to the channels enum
ALTER TYPE channels ADD VALUE IF NOT EXISTS 'whatsapp';
ALTER TYPE channels ADD VALUE IF NOT EXISTS 'telegram';
ALTER TYPE channels ADD VALUE IF NOT EXISTS 'sms';
ALTER TYPE channels ADD VALUE IF NOT EXISTS 'push';
