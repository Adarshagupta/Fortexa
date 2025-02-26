# Fortexa Merchant Dashboard

A modern React-based dashboard for the Fortexa payment processing platform. This application provides merchants with a comprehensive interface to manage payments, settlements, and account settings.

## Features

- **Secure Authentication**: User login with JWT-based authentication
- **Dashboard Overview**: Real-time metrics and charts showing payment activity
- **Payment Management**: View, filter, and manage payment transactions
- **Settlement Tracking**: Monitor settlement status and financial reconciliation
- **Responsive Design**: Fully responsive UI that works on desktop and mobile devices
- **Data Visualization**: Interactive charts for better data understanding

## Tech Stack

- **Frontend**: React, TypeScript, Material-UI
- **State Management**: React Hooks (Context API)
- **Routing**: React Router v6
- **Styling**: Material-UI with custom theming
- **Charts**: Recharts
- **HTTP Client**: Axios

## Project Structure

```
src/
├── components/         # Reusable UI components
│   ├── Layout/         # App layout components (Header, Sidebar, etc.)
│   └── common/         # Common UI elements
├── pages/              # Top-level page components
│   ├── Dashboard/      # Dashboard page
│   ├── Payments/       # Payments management page
│   ├── Settlements/    # Settlements management page
│   └── Login/          # Authentication page
├── services/           # API services
├── types/              # TypeScript interfaces and types
├── theme/              # Material-UI theme configuration
├── utils/              # Utility functions
└── App.tsx             # Main application component
```

## Getting Started

### Prerequisites

- Node.js (v16 or later)
- npm or yarn

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/fortexa-merchant-dashboard.git
cd fortexa-merchant-dashboard
```

2. Install dependencies
```bash
npm install
# or
yarn install
```

3. Start the development server
```bash
npm start
# or
yarn start
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser

### Demo Credentials

For testing purposes, you can use the following credentials:
- Email: `demo@fortexa.com`
- Password: `password123`

## API Integration

The dashboard connects to the Fortexa Payment API. In development mode, it uses mock data for demonstration purposes. In production, it will connect to the actual API endpoints.

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
REACT_APP_API_URL=https://api.fortexa.example.com
```

## Deployment

To build the application for production:

```bash
npm run build
# or
yarn build
```

This will create a `build` folder with optimized production build.

## Development Guidelines

- Follow the established component structure
- Add proper TypeScript types for all components and functions
- Use Material-UI components and theme variables for consistent styling
- Write unit tests for critical functionality
- Document complex logic with comments

## License

This project is licensed under the MIT License - see the LICENSE file for details.
