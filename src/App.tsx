import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, CssBaseline } from '@mui/material';

// Theme
import theme from './theme';

// Layout
import { Layout } from './components/Layout';

// Pages
import Dashboard from './pages/Dashboard';
import Payments from './pages/Payments';
import Settlements from './pages/Settlements';
import Login from './pages/Login';

// Auth guard - checks if user is authenticated
const AuthGuard = ({ children }: { children: React.ReactNode }) => {
  // In a real app, check for authentication
  const isAuthenticated = localStorage.getItem('apiKey') !== null;
  
  if (!isAuthenticated) {
    // Redirect to login if not authenticated
    return <Navigate to="/login" />;
  }
  
  return <>{children}</>;
};

// Guard to prevent authenticated users from accessing login page
const NoAuthGuard = ({ children }: { children: React.ReactNode }) => {
  const isAuthenticated = localStorage.getItem('apiKey') !== null;
  
  if (isAuthenticated) {
    // Redirect to dashboard if already authenticated
    return <Navigate to="/" />;
  }
  
  return <>{children}</>;
};

// Main App component
const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <Routes>
          {/* Public routes */}
          <Route
            path="/login"
            element={
              <NoAuthGuard>
                <Login />
              </NoAuthGuard>
            }
          />

          {/* Protected routes that require authentication */}
          <Route
            path="/"
            element={
              <AuthGuard>
                <Layout>
                  <Dashboard />
                </Layout>
              </AuthGuard>
            }
          />
          <Route
            path="/payments"
            element={
              <AuthGuard>
                <Layout>
                  <Payments />
                </Layout>
              </AuthGuard>
            }
          />
          <Route
            path="/settlements"
            element={
              <AuthGuard>
                <Layout>
                  <Settlements />
                </Layout>
              </AuthGuard>
            }
          />
          
          {/* Add more routes as needed */}
          <Route
            path="*"
            element={<Navigate to="/" replace />}
          />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
};

export default App; 