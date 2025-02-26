import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Chip,
  TextField,
  MenuItem,
  Button,
  Grid,
  IconButton,
  Tooltip,
  CircularProgress,
} from '@mui/material';
import {
  Search as SearchIcon,
  FilterList as FilterIcon,
  Refresh as RefreshIcon,
  MoreVert as MoreVertIcon,
  Visibility as VisibilityIcon,
  GetApp as DownloadIcon,
} from '@mui/icons-material';

import { paymentsApi } from '../../services/api';
import { Payment, PaymentStatus, PaymentMethod } from '../../types';

// Styles for payment status chips
const statusColorMap: Record<PaymentStatus, string> = {
  [PaymentStatus.INITIATED]: 'info',
  [PaymentStatus.AUTHORIZED]: 'warning',
  [PaymentStatus.CAPTURED]: 'success',
  [PaymentStatus.SETTLED]: 'success',
  [PaymentStatus.REFUNDED]: 'secondary',
  [PaymentStatus.FAILED]: 'error',
  [PaymentStatus.CHARGEBACK]: 'error',
};

const Payments: React.FC = () => {
  // State
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [totalCount, setTotalCount] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<PaymentStatus | 'ALL'>('ALL');

  // Fetch payments
  useEffect(() => {
    const fetchPayments = async () => {
      setLoading(true);
      try {
        const filters: Record<string, any> = {};
        if (statusFilter !== 'ALL') {
          filters.status = statusFilter;
        }
        if (searchQuery) {
          filters.search = searchQuery;
        }

        // Mock data for demo purposes
        // In a real app, we would use:
        // const response = await paymentsApi.getPayments(page + 1, rowsPerPage, filters);
        // if (response.success) {
        //   setPayments(response.data.payments);
        //   setTotalCount(response.data.total);
        // }

        // Mock payments data
        const mockPayments: Payment[] = [
          {
            id: '1',
            merchantId: 'merchant-1',
            customerId: 'customer-1',
            amount: 1250.00,
            currency: 'USD',
            status: PaymentStatus.CAPTURED,
            paymentMethodId: 'pm-1',
            paymentMethodType: PaymentMethod.CREDIT_CARD,
            description: 'Subscription payment',
            referenceId: 'ref-1',
            createdAt: '2023-01-10T10:30:00Z',
            updatedAt: '2023-01-10T10:35:00Z',
          },
          {
            id: '2',
            merchantId: 'merchant-1',
            customerId: 'customer-2',
            amount: 890.50,
            currency: 'USD',
            status: PaymentStatus.SETTLED,
            paymentMethodId: 'pm-2',
            paymentMethodType: PaymentMethod.BANK_TRANSFER,
            description: 'One-time purchase',
            referenceId: 'ref-2',
            createdAt: '2023-01-09T09:15:00Z',
            updatedAt: '2023-01-09T09:20:00Z',
          },
          {
            id: '3',
            merchantId: 'merchant-1',
            customerId: 'customer-3',
            amount: 450.00,
            currency: 'USD',
            status: PaymentStatus.FAILED,
            paymentMethodId: 'pm-3',
            paymentMethodType: PaymentMethod.UPI,
            description: 'Service payment',
            referenceId: 'ref-3',
            createdAt: '2023-01-08T18:45:00Z',
            updatedAt: '2023-01-08T18:50:00Z',
          },
          {
            id: '4',
            merchantId: 'merchant-1',
            customerId: 'customer-4',
            amount: 2100.75,
            currency: 'USD',
            status: PaymentStatus.AUTHORIZED,
            paymentMethodId: 'pm-4',
            paymentMethodType: PaymentMethod.DEBIT_CARD,
            description: 'Product purchase',
            referenceId: 'ref-4',
            createdAt: '2023-01-07T15:20:00Z',
            updatedAt: '2023-01-07T15:25:00Z',
          },
          {
            id: '5',
            merchantId: 'merchant-1',
            customerId: 'customer-5',
            amount: 75.25,
            currency: 'USD',
            status: PaymentStatus.INITIATED,
            paymentMethodId: 'pm-5',
            paymentMethodType: PaymentMethod.CREDIT_CARD,
            description: 'Subscription renewal',
            referenceId: 'ref-5',
            createdAt: '2023-01-06T12:10:00Z',
            updatedAt: '2023-01-06T12:15:00Z',
          },
        ];

        setPayments(mockPayments);
        setTotalCount(50); // Mock total count
        setLoading(false);
      } catch (error) {
        console.error('Error fetching payments:', error);
        setLoading(false);
      }
    };

    fetchPayments();
  }, [page, rowsPerPage, searchQuery, statusFilter]);

  // Handle page change
  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  // Handle rows per page change
  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  // Format currency
  const formatCurrency = (amount: number, currency: string): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(amount);
  };

  // Format date
  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Payments
      </Typography>

      {/* Filters Section */}
      <Paper
        elevation={0}
        sx={{
          p: 2,
          mb: 3,
          borderRadius: 2,
          boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
        }}
      >
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              size="small"
              label="Search"
              variant="outlined"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                startAdornment: <SearchIcon fontSize="small" sx={{ mr: 1, color: 'text.secondary' }} />,
              }}
              placeholder="Search by ID, customer, reference..."
            />
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              select
              fullWidth
              size="small"
              label="Payment Status"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as PaymentStatus | 'ALL')}
              InputProps={{
                startAdornment: <FilterIcon fontSize="small" sx={{ mr: 1, color: 'text.secondary' }} />,
              }}
            >
              <MenuItem value="ALL">All Statuses</MenuItem>
              {Object.values(PaymentStatus).map((status) => (
                <MenuItem key={status} value={status}>
                  {status}
                </MenuItem>
              ))}
            </TextField>
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              select
              fullWidth
              size="small"
              label="Date Range"
              value="LAST_30_DAYS"
            >
              <MenuItem value="TODAY">Today</MenuItem>
              <MenuItem value="YESTERDAY">Yesterday</MenuItem>
              <MenuItem value="LAST_7_DAYS">Last 7 days</MenuItem>
              <MenuItem value="LAST_30_DAYS">Last 30 days</MenuItem>
              <MenuItem value="CUSTOM">Custom range</MenuItem>
            </TextField>
          </Grid>
          <Grid item xs={12} md={2} sx={{ display: 'flex', justifyContent: { xs: 'flex-start', md: 'flex-end' } }}>
            <Button
              variant="outlined"
              color="primary"
              startIcon={<RefreshIcon />}
              onClick={() => {
                setSearchQuery('');
                setStatusFilter('ALL');
                setPage(0);
              }}
            >
              Reset
            </Button>
          </Grid>
        </Grid>
      </Paper>

      {/* Payments Table */}
      <Paper
        elevation={0}
        sx={{
          borderRadius: 2,
          overflow: 'hidden',
          boxShadow: '0 2px 10px rgba(0, 0, 0, 0.05)',
        }}
      >
        <TableContainer>
          <Table>
            <TableHead sx={{ backgroundColor: 'primary.light', color: 'primary.contrastText' }}>
              <TableRow>
                <TableCell sx={{ fontWeight: 'bold' }}>Payment ID</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Date</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Customer</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Amount</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Payment Method</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Status</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 3 }}>
                    <CircularProgress size={40} />
                  </TableCell>
                </TableRow>
              ) : payments.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 3 }}>
                    <Typography variant="body1">No payments found</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                payments.map((payment) => (
                  <TableRow key={payment.id} hover>
                    <TableCell>{payment.id}</TableCell>
                    <TableCell>{formatDate(payment.createdAt)}</TableCell>
                    <TableCell>{payment.customerId || 'N/A'}</TableCell>
                    <TableCell>{formatCurrency(payment.amount, payment.currency)}</TableCell>
                    <TableCell>
                      {payment.paymentMethodType.replace('_', ' ')}
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={payment.status}
                        size="small"
                        color={statusColorMap[payment.status] as any}
                        sx={{ fontWeight: 500 }}
                      />
                    </TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex' }}>
                        <Tooltip title="View Details">
                          <IconButton size="small" color="primary">
                            <VisibilityIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Download Receipt">
                          <IconButton size="small" color="primary">
                            <DownloadIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="More Actions">
                          <IconButton size="small">
                            <MoreVertIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                      </Box>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
        
        <TablePagination
          rowsPerPageOptions={[5, 10, 25, 50]}
          component="div"
          count={totalCount}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </Paper>
    </Box>
  );
};

export default Payments; 