// Payment related types
export enum PaymentStatus {
  INITIATED = "INITIATED",
  AUTHORIZED = "AUTHORIZED",
  CAPTURED = "CAPTURED",
  SETTLED = "SETTLED",
  REFUNDED = "REFUNDED",
  FAILED = "FAILED",
  CHARGEBACK = "CHARGEBACK"
}

export enum PaymentMethod {
  CREDIT_CARD = "CREDIT_CARD",
  DEBIT_CARD = "DEBIT_CARD",
  UPI = "UPI",
  BANK_TRANSFER = "BANK_TRANSFER",
  WALLET = "WALLET",
  CRYPTO = "CRYPTO",
  BNPL = "BNPL"
}

export interface Payment {
  id: string;
  merchantId: string;
  customerId?: string;
  amount: number;
  currency: string;
  status: PaymentStatus;
  paymentMethodId?: string;
  paymentMethodType: PaymentMethod;
  description?: string;
  metadata?: Record<string, any>;
  idempotencyKey?: string;
  referenceId?: string;
  createdAt: string;
  updatedAt: string;
}

// Settlement related types
export enum SettlementStatus {
  PENDING = "PENDING",
  PROCESSING = "PROCESSING",
  COMPLETED = "COMPLETED",
  FAILED = "FAILED"
}

export enum SettlementMethod {
  BANK_TRANSFER = "BANK_TRANSFER",
  WALLET = "WALLET"
}

export interface Settlement {
  id: string;
  merchantId: string;
  amount: number;
  currency: string;
  status: SettlementStatus;
  paymentCount: number;
  feeAmount: number;
  taxAmount: number;
  netAmount: number;
  settlementDate: string;
  bankAccountId: string;
  settlementMethod: SettlementMethod;
  reference: string;
  createdAt: string;
  updatedAt: string;
}

// Merchant related types
export enum MerchantStatus {
  ACTIVE = "ACTIVE",
  INACTIVE = "INACTIVE",
  PENDING = "PENDING",
  BLOCKED = "BLOCKED"
}

export interface Merchant {
  id: string;
  name: string;
  businessName: string;
  email: string;
  phone?: string;
  website?: string;
  apiKey: string;
  status: MerchantStatus;
  createdAt: string;
  updatedAt: string;
}

// Dashboard stats
export interface DashboardStats {
  totalPayments: number;
  totalAmount: number;
  successRate: number;
  pendingSettlements: number;
  recentPayments: Payment[];
  paymentsByMethod: Record<PaymentMethod, number>;
  revenueByDay: Array<{
    date: string;
    amount: number;
  }>;
}

// API response type
export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
  errors?: any;
} 