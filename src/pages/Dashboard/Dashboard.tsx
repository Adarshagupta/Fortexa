import React, { useEffect, useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Grid,
  Typography,
  CircularProgress,
  Button,
  Paper,
  Divider,
  IconButton,
} from '@mui/material';
import {
  MoreVert as MoreVertIcon,
  TrendingUp as TrendingUpIcon,
  AttachMoney as MoneyIcon,
  Payment as PaymentIcon,
  AccountBalanceWallet as WalletIcon,
} from '@mui/icons-material';
import { Link as RouterLink } from 'react-router-dom';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

import { dashboardApi } from '../../services/api';
import { DashboardStats, Payment, PaymentMethod, PaymentStatus } from '../../types';

// Mock data for testing
const mockRevenueData = [
  { date: '2023-01-01', amount: 1500 },
  { date: '2023-01-02', amount: 2300 },
  { date: '2023-01-03', amount: 1800 },
  { date: '2023-01-04', amount: 2800 },
  { date: '2023-01-05', amount: 3200 },
  { date: '2023-01-06', amount: 2100 },
  { date: '2023-01-07', amount: 3500 },
];

const mockPaymentMethodsData = [
  { name: 'Credit Card', value: 45 },
  { name: 'Debit Card', value: 25 },
  { name: 'UPI', value: 20 },
  { name: 'Bank Transfer', value: 10 },
];

const mockRecentPayments: Payment[] = [
  {
    id: '1',
    merchantId: 'merchant-1',
    amount: 1250.00,
    currency: 'USD',
    status: PaymentStatus.CAPTURED,
    paymentMethodType: PaymentMethod.CREDIT_CARD,
    createdAt: '2023-01-07T10:30:00Z',
    updatedAt: '2023-01-07T10:35:00Z',
  },
  {
    id: '2',
    merchantId: 'merchant-1',
    amount: 890.50,
    currency: 'USD',
    status: PaymentStatus.SETTLED,
    paymentMethodType: PaymentMethod.BANK_TRANSFER,
    createdAt: '2023-01-07T09:15:00Z',
    updatedAt: '2023-01-07T09:20:00Z',
  },
  {
    id: '3',
    merchantId: 'merchant-1',
    amount: 450.00,
    currency: 'USD',
    status: PaymentStatus.FAILED,
    paymentMethodType: PaymentMethod.UPI,
    createdAt: '2023-01-06T18:45:00Z',
    updatedAt: '2023-01-06T18:50:00Z',
  },
  {
    id: '4',
    merchantId: 'merchant-1',
    amount: 2100.75,
    currency: 'USD',
    status: PaymentStatus.AUTHORIZED,
    paymentMethodType: PaymentMethod.DEBIT_CARD,
    createdAt: '2023-01-06T15:20:00Z',
    updatedAt: '2023-01-06T15:25:00Z',
  },
];

// Colors for the pie chart
const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<DashboardStats | null>(null);
  
  // Fetch dashboard stats
  useEffect(() => {
    const fetchStats = async () => {
      try {
        // For now, we'll use mock data instead of making actual API calls
        // const response = await dashboardApi.getStats();
        // if (response.success) {
        //   setStats(response.data);
        // }
        
        // Mock data
        setStats({
          totalPayments: 247,
          totalAmount: 53850.75,
          successRate: 94.5,
          pendingSettlements: 2,
          recentPayments: mockRecentPayments,
          paymentsByMethod: {
            [PaymentMethod.CREDIT_CARD]: 120,
            [PaymentMethod.DEBIT_CARD]: 65,
            [PaymentMethod.UPI]: 40,
            [PaymentMethod.BANK_TRANSFER]: 22,
            [PaymentMethod.WALLET]: 0,
            [PaymentMethod.CRYPTO]: 0,
            [PaymentMethod.BNPL]: 0,
          },
          revenueByDay: mockRevenueData,
        });
        
        setLoading(false);
      } catch (error) {
        console.error('Error fetching dashboard stats:', error);
        setLoading(false);
      }
    };
    
    fetchStats();
  }, []);
  
  // Format currency
  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };
  
  // Format date
  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    });
  };
  
  if (loading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100%',
        }}
      >
        <CircularProgress />
      </Box>
    );
  }
  
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      
      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card 
            elevation={0}
            sx={{ 
              borderRadius: 2, 
              height: '100%',
              boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-5px)',
                boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
              }
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Total Transactions
                </Typography>
                <MoneyIcon color="primary" />
              </Box>
              <Typography variant="h4" sx={{ fontWeight: 600 }}>
                {stats?.totalPayments || 0}
              </Typography>
              <Typography variant="body2" sx={{ mt: 1, color: 'success.main' }}>
                <TrendingUpIcon fontSize="small" sx={{ verticalAlign: 'middle', mr: 0.5 }} />
                +12.5% from last week
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} sm={6} md={3}>
          <Card 
            elevation={0}
            sx={{ 
              borderRadius: 2, 
              height: '100%',
              boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-5px)',
                boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
              }
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Total Revenue
                </Typography>
                <PaymentIcon color="primary" />
              </Box>
              <Typography variant="h4" sx={{ fontWeight: 600 }}>
                {formatCurrency(stats?.totalAmount || 0)}
              </Typography>
              <Typography variant="body2" sx={{ mt: 1, color: 'success.main' }}>
                <TrendingUpIcon fontSize="small" sx={{ verticalAlign: 'middle', mr: 0.5 }} />
                +8.2% from last week
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} sm={6} md={3}>
          <Card 
            elevation={0}
            sx={{ 
              borderRadius: 2, 
              height: '100%',
              boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-5px)',
                boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
              }
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Success Rate
                </Typography>
                <WalletIcon color="primary" />
              </Box>
              <Typography variant="h4" sx={{ fontWeight: 600 }}>
                {stats?.successRate || 0}%
              </Typography>
              <Typography variant="body2" sx={{ mt: 1, color: 'success.main' }}>
                <TrendingUpIcon fontSize="small" sx={{ verticalAlign: 'middle', mr: 0.5 }} />
                +1.5% from last week
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} sm={6} md={3}>
          <Card 
            elevation={0}
            sx={{ 
              borderRadius: 2, 
              height: '100%',
              boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-5px)',
                boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
              }
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Pending Settlements
                </Typography>
                <WalletIcon color="primary" />
              </Box>
              <Typography variant="h4" sx={{ fontWeight: 600 }}>
                {stats?.pendingSettlements || 0}
              </Typography>
              <Button
                component={RouterLink}
                to="/settlements"
                variant="text"
                color="primary"
                size="small"
                sx={{ mt: 1 }}
              >
                View Details
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      {/* Charts */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={8}>
          <Paper 
            elevation={0}
            sx={{ 
              p: 3,
              height: '100%',
              borderRadius: 2,
              boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
            }}
          >
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
              <Typography variant="h6">Revenue Trend</Typography>
              <IconButton size="small">
                <MoreVertIcon fontSize="small" />
              </IconButton>
            </Box>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart
                data={stats?.revenueByDay || mockRevenueData}
                margin={{
                  top: 5,
                  right: 30,
                  left: 20,
                  bottom: 5,
                }}
              >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="date" 
                  tickFormatter={(dateStr) => new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} 
                />
                <YAxis />
                <Tooltip 
                  formatter={(value) => [formatCurrency(Number(value)), 'Revenue']} 
                  labelFormatter={(dateStr) => new Date(dateStr).toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })} 
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="amount"
                  name="Revenue"
                  stroke="#1976d2"
                  strokeWidth={2}
                  activeDot={{ r: 8 }}
                />
              </LineChart>
            </ResponsiveContainer>
          </Paper>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Paper 
            elevation={0}
            sx={{ 
              p: 3,
              height: '100%',
              borderRadius: 2,
              boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
            }}
          >
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
              <Typography variant="h6">Payment Methods</Typography>
              <IconButton size="small">
                <MoreVertIcon fontSize="small" />
              </IconButton>
            </Box>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={mockPaymentMethodsData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={90}
                  fill="#8884d8"
                  paddingAngle={5}
                  dataKey="value"
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                >
                  {mockPaymentMethodsData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [`${value}%`, 'Percentage']} />
              </PieChart>
            </ResponsiveContainer>
          </Paper>
        </Grid>
      </Grid>
      
      {/* Recent Transactions */}
      <Paper 
        elevation={0}
        sx={{ 
          p: 3,
          borderRadius: 2,
          boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
        }}
      >
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h6">Recent Transactions</Typography>
          <Button 
            component={RouterLink}
            to="/payments"
            variant="outlined" 
            color="primary" 
            size="small"
          >
            View All
          </Button>
        </Box>
        
        {stats?.recentPayments.map((payment, index) => (
          <React.Fragment key={payment.id}>
            <Box sx={{ py: 1.5, display: 'flex', justifyContent: 'space-between' }}>
              <Box>
                <Typography variant="subtitle2">
                  Payment #{payment.id}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  {formatDate(payment.createdAt)} via {payment.paymentMethodType.replace('_', ' ')}
                </Typography>
              </Box>
              <Box sx={{ textAlign: 'right' }}>
                <Typography variant="subtitle2">
                  {formatCurrency(payment.amount)}
                </Typography>
                <Typography 
                  variant="body2" 
                  sx={{ 
                    color: 
                      payment.status === PaymentStatus.CAPTURED || payment.status === PaymentStatus.SETTLED
                        ? 'success.main'
                        : payment.status === PaymentStatus.FAILED
                        ? 'error.main'
                        : 'warning.main'
                  }}
                >
                  {payment.status}
                </Typography>
              </Box>
            </Box>
            {index < stats.recentPayments.length - 1 && <Divider />}
          </React.Fragment>
        ))}
      </Paper>
    </Box>
  );
};

export default Dashboard; 