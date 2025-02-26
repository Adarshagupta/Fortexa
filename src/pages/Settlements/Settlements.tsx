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
  Card,
  CardContent,
} from '@mui/material';
import {
  Search as SearchIcon,
  FilterList as FilterIcon,
  Refresh as RefreshIcon,
  MoreVert as MoreVertIcon,
  Visibility as VisibilityIcon,
  GetApp as DownloadIcon,
  AccountBalance as AccountBalanceIcon,
  CalendarToday as CalendarIcon,
  AttachMoney as MoneyIcon,
} from '@mui/icons-material';

import { settlementsApi } from '../../services/api';
import { Settlement, SettlementStatus, SettlementMethod } from '../../types';

// Styles for settlement status chips
const statusColorMap: Record<SettlementStatus, string> = {
  [SettlementStatus.PENDING]: 'warning',
  [SettlementStatus.PROCESSING]: 'info',
  [SettlementStatus.COMPLETED]: 'success',
  [SettlementStatus.FAILED]: 'error',
};

const Settlements: React.FC = () => {
  // State
  const [settlements, setSettlements] = useState<Settlement[]>([]);
  const [loading, setLoading] = useState(true);
  const [totalCount, setTotalCount] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<SettlementStatus | 'ALL'>('ALL');

  // Summary stats
  const [summaryStats, setSummaryStats] = useState({
    totalAmount: 0,
    completedAmount: 0,
    pendingAmount: 0,
    settlementCount: 0,
  });

  // Fetch settlements
  useEffect(() => {
    const fetchSettlements = async () => {
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
        // const response = await settlementsApi.getSettlements(page + 1, rowsPerPage, filters);
        // if (response.success) {
        //   setSettlements(response.data.settlements);
        //   setTotalCount(response.data.total);
        // }

        // Mock settlements data
        const mockSettlements: Settlement[] = [
          {
            id: 'stl-1',
            merchantId: 'merchant-1',
            amount: 5250.00,
            currency: 'USD',
            status: SettlementStatus.COMPLETED,
            paymentCount: 12,
            feeAmount: 105.00,
            taxAmount: 0,
            netAmount: 5145.00,
            settlementDate: '2023-01-15T10:30:00Z',
            bankAccountId: 'ba-1',
            settlementMethod: SettlementMethod.BANK_TRANSFER,
            reference: 'ref-stl-1',
            createdAt: '2023-01-15T10:30:00Z',
            updatedAt: '2023-01-15T10:35:00Z',
          },
          {
            id: 'stl-2',
            merchantId: 'merchant-1',
            amount: 3890.50,
            currency: 'USD',
            status: SettlementStatus.PENDING,
            paymentCount: 8,
            feeAmount: 77.81,
            taxAmount: 0,
            netAmount: 3812.69,
            settlementDate: '2023-01-20T00:00:00Z',
            bankAccountId: 'ba-1',
            settlementMethod: SettlementMethod.BANK_TRANSFER,
            reference: 'ref-stl-2',
            createdAt: '2023-01-18T09:15:00Z',
            updatedAt: '2023-01-18T09:20:00Z',
          },
          {
            id: 'stl-3',
            merchantId: 'merchant-1',
            amount: 2450.00,
            currency: 'USD',
            status: SettlementStatus.PROCESSING,
            paymentCount: 5,
            feeAmount: 49.00,
            taxAmount: 0,
            netAmount: 2401.00,
            settlementDate: '2023-01-17T00:00:00Z',
            bankAccountId: 'ba-1',
            settlementMethod: SettlementMethod.BANK_TRANSFER,
            reference: 'ref-stl-3',
            createdAt: '2023-01-16T18:45:00Z',
            updatedAt: '2023-01-16T18:50:00Z',
          },
          {
            id: 'stl-4',
            merchantId: 'merchant-1',
            amount: 1200.75,
            currency: 'USD',
            status: SettlementStatus.COMPLETED,
            paymentCount: 3,
            feeAmount: 24.01,
            taxAmount: 0,
            netAmount: 1176.74,
            settlementDate: '2023-01-10T00:00:00Z',
            bankAccountId: 'ba-1',
            settlementMethod: SettlementMethod.BANK_TRANSFER,
            reference: 'ref-stl-4',
            createdAt: '2023-01-09T15:20:00Z',
            updatedAt: '2023-01-10T15:25:00Z',
          },
          {
            id: 'stl-5',
            merchantId: 'merchant-1',
            amount: 8750.25,
            currency: 'USD',
            status: SettlementStatus.FAILED,
            paymentCount: 20,
            feeAmount: 175.01,
            taxAmount: 0,
            netAmount: 8575.24,
            settlementDate: '2023-01-05T00:00:00Z',
            bankAccountId: 'ba-1',
            settlementMethod: SettlementMethod.BANK_TRANSFER,
            reference: 'ref-stl-5',
            createdAt: '2023-01-04T12:10:00Z',
            updatedAt: '2023-01-05T12:15:00Z',
          },
        ];

        setSettlements(mockSettlements);
        setTotalCount(20); // Mock total count

        // Calculate summary stats
        const totalAmount = mockSettlements.reduce((sum, item) => sum + item.amount, 0);
        const completedAmount = mockSettlements
          .filter(item => item.status === SettlementStatus.COMPLETED)
          .reduce((sum, item) => sum + item.amount, 0);
        const pendingAmount = mockSettlements
          .filter(item => item.status === SettlementStatus.PENDING || item.status === SettlementStatus.PROCESSING)
          .reduce((sum, item) => sum + item.amount, 0);

        setSummaryStats({
          totalAmount,
          completedAmount,
          pendingAmount,
          settlementCount: mockSettlements.length,
        });

        setLoading(false);
      } catch (error) {
        console.error('Error fetching settlements:', error);
        setLoading(false);
      }
    };

    fetchSettlements();
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
    });
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Settlements
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
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Total Volume
                </Typography>
                <MoneyIcon color="primary" />
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600 }}>
                {formatCurrency(summaryStats.totalAmount, 'USD')}
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
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Completed
                </Typography>
                <AccountBalanceIcon color="success" />
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600 }}>
                {formatCurrency(summaryStats.completedAmount, 'USD')}
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
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Pending
                </Typography>
                <CalendarIcon color="warning" />
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600 }}>
                {formatCurrency(summaryStats.pendingAmount, 'USD')}
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
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography color="textSecondary" variant="subtitle2">
                  Total Settlements
                </Typography>
                <AccountBalanceIcon color="primary" />
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600 }}>
                {summaryStats.settlementCount}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

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
              placeholder="Search by ID, reference..."
            />
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              select
              fullWidth
              size="small"
              label="Settlement Status"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as SettlementStatus | 'ALL')}
              InputProps={{
                startAdornment: <FilterIcon fontSize="small" sx={{ mr: 1, color: 'text.secondary' }} />,
              }}
            >
              <MenuItem value="ALL">All Statuses</MenuItem>
              {Object.values(SettlementStatus).map((status) => (
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

      {/* Settlements Table */}
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
                <TableCell sx={{ fontWeight: 'bold' }}>Settlement ID</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Date</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Amount</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Net Amount</TableCell>
                <TableCell sx={{ fontWeight: 'bold' }}>Payment Count</TableCell>
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
              ) : settlements.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 3 }}>
                    <Typography variant="body1">No settlements found</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                settlements.map((settlement) => (
                  <TableRow key={settlement.id} hover>
                    <TableCell>{settlement.id}</TableCell>
                    <TableCell>{formatDate(settlement.settlementDate)}</TableCell>
                    <TableCell>{formatCurrency(settlement.amount, settlement.currency)}</TableCell>
                    <TableCell>{formatCurrency(settlement.netAmount, settlement.currency)}</TableCell>
                    <TableCell>{settlement.paymentCount}</TableCell>
                    <TableCell>
                      <Chip
                        label={settlement.status}
                        size="small"
                        color={statusColorMap[settlement.status] as any}
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
                        <Tooltip title="Download Report">
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

export default Settlements; 