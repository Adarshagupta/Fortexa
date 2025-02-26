-- Initialize Fortexa Database Schema

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types for status tracking
CREATE TYPE payment_status AS ENUM (
  'INITIATED',
  'AUTHORIZED',
  'CAPTURED',
  'SETTLED',
  'REFUNDED',
  'FAILED',
  'CHARGEBACK'
);

CREATE TYPE payment_method AS ENUM (
  'CREDIT_CARD',
  'DEBIT_CARD',
  'UPI',
  'BANK_TRANSFER',
  'WALLET',
  'CRYPTO',
  'BNPL'
);

CREATE TYPE settlement_status AS ENUM (
  'PENDING',
  'PROCESSING',
  'COMPLETED',
  'FAILED'
);

-- Create merchants table
CREATE TABLE merchants (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(100) NOT NULL,
  business_name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  phone VARCHAR(20),
  website VARCHAR(100),
  api_key VARCHAR(64) UNIQUE NOT NULL,
  secret_key VARCHAR(64) NOT NULL,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create customers table
CREATE TABLE customers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(100),
  phone VARCHAR(20),
  name VARCHAR(100),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT email_or_phone CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

-- Create payment_methods table
CREATE TABLE payment_methods (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  customer_id UUID REFERENCES customers(id),
  type payment_method NOT NULL,
  token VARCHAR(100),
  card_last_four VARCHAR(4),
  card_expiry_month VARCHAR(2),
  card_expiry_year VARCHAR(4),
  card_brand VARCHAR(20),
  upi_id VARCHAR(50),
  bank_account_number VARCHAR(50),
  bank_ifsc VARCHAR(20),
  wallet_id VARCHAR(50),
  crypto_address VARCHAR(100),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create payments table
CREATE TABLE payments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  merchant_id UUID REFERENCES merchants(id) NOT NULL,
  customer_id UUID REFERENCES customers(id),
  amount DECIMAL(12, 2) NOT NULL,
  currency VARCHAR(3) NOT NULL DEFAULT 'INR',
  status payment_status NOT NULL DEFAULT 'INITIATED',
  payment_method_id UUID REFERENCES payment_methods(id),
  payment_method_type payment_method NOT NULL,
  description TEXT,
  metadata JSONB,
  idempotency_key VARCHAR(100) UNIQUE,
  reference_id VARCHAR(100),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create transactions table to track state changes
CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  payment_id UUID REFERENCES payments(id) NOT NULL,
  amount DECIMAL(12, 2) NOT NULL,
  status VARCHAR(20) NOT NULL,
  gateway_response JSONB,
  error_message TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create settlements table
CREATE TABLE settlements (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  merchant_id UUID REFERENCES merchants(id) NOT NULL,
  amount DECIMAL(12, 2) NOT NULL,
  fees DECIMAL(12, 2) NOT NULL DEFAULT 0,
  net_amount DECIMAL(12, 2) NOT NULL,
  status settlement_status NOT NULL DEFAULT 'PENDING',
  settlement_date TIMESTAMP WITH TIME ZONE,
  transaction_count INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create settlement_items table
CREATE TABLE settlement_items (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  settlement_id UUID REFERENCES settlements(id) NOT NULL,
  payment_id UUID REFERENCES payments(id) NOT NULL,
  amount DECIMAL(12, 2) NOT NULL,
  fees DECIMAL(12, 2) NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create webhooks table for merchant notifications
CREATE TABLE webhooks (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  merchant_id UUID REFERENCES merchants(id) NOT NULL,
  url VARCHAR(255) NOT NULL,
  event_types VARCHAR[] NOT NULL,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  secret VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create webhook_events table for tracking webhook deliveries
CREATE TABLE webhook_events (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  webhook_id UUID REFERENCES webhooks(id) NOT NULL,
  event_type VARCHAR(50) NOT NULL,
  payload JSONB NOT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  attempts INTEGER DEFAULT 0,
  last_attempt_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_payments_merchant_id ON payments(merchant_id);
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_transactions_payment_id ON transactions(payment_id);
CREATE INDEX idx_settlements_merchant_id ON settlements(merchant_id);
CREATE INDEX idx_settlement_items_settlement_id ON settlement_items(settlement_id);
CREATE INDEX idx_webhook_events_webhook_id ON webhook_events(webhook_id);

-- Insert sample merchant for testing
INSERT INTO merchants (name, business_name, email, phone, website, api_key, secret_key)
VALUES (
  'Test Merchant',
  'Test Business',
  'test@example.com',
  '+919999999999',
  'https://example.com',
  'pk_test_' || md5(random()::text),
  'sk_test_' || md5(random()::text)
);

COMMENT ON TABLE merchants IS 'Stores merchant information';
COMMENT ON TABLE customers IS 'Stores customer information';
COMMENT ON TABLE payment_methods IS 'Stores customer payment methods';
COMMENT ON TABLE payments IS 'Stores payment transactions';
COMMENT ON TABLE transactions IS 'Stores transaction state changes';
COMMENT ON TABLE settlements IS 'Stores merchant settlements';
COMMENT ON TABLE settlement_items IS 'Stores individual items in a settlement';
COMMENT ON TABLE webhooks IS 'Stores merchant webhook configurations';
COMMENT ON TABLE webhook_events IS 'Stores webhook event delivery attempts'; 