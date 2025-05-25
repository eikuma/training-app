import { createTheme } from '@mui/material/styles';
import { red } from '@mui/material/colors'; // For error states

// Define a clean, modern font stack
const FONT_FAMILY = [
  '-apple-system',
  'BlinkMacSystemFont',
  '"Segoe UI"',
  'Roboto',
  '"Helvetica Neue"',
  'Arial',
  'sans-serif',
  '"Apple Color Emoji"',
  '"Segoe UI Emoji"',
  '"Segoe UI Symbol"',
].join(',');

const theme = createTheme({
  palette: {
    primary: {
      main: '#00796b', // A sophisticated teal
      // light: will be calculated from main,
      // dark: will be calculated from main,
      // contrastText: will be calculated to contrast with main
    },
    secondary: {
      main: '#fbc02d', // A vibrant amber/yellow
    },
    error: {
      main: red.A400, // Standard Material Design error color
    },
    background: {
      default: '#f4f6f8', // A very light grey, slightly off-white
      paper: '#ffffff',   // Standard white for paper elements
    },
    text: {
      primary: '#333333', // Dark grey for primary text for better readability
      secondary: '#575757', // Slightly lighter grey for secondary text
    },
  },
  typography: {
    fontFamily: FONT_FAMILY,
    h1: {
      fontSize: '2.75rem',
      fontWeight: 600,
      lineHeight: 1.2,
      letterSpacing: '-0.01562em',
    },
    h2: {
      fontSize: '2.25rem',
      fontWeight: 600,
      lineHeight: 1.25,
      letterSpacing: '-0.00833em',
    },
    h3: {
      fontSize: '1.75rem',
      fontWeight: 600,
      lineHeight: 1.3,
      letterSpacing: '0em',
    },
    h4: {
      fontSize: '1.5rem', // Used for page titles in recommendations/page.tsx
      fontWeight: 500, // Slightly less bold than h1-h3
      lineHeight: 1.33,
      letterSpacing: '0.00735em',
    },
    h5: {
      fontSize: '1.25rem', // Used in workouts/page.tsx for titles
      fontWeight: 500,
      lineHeight: 1.35,
      letterSpacing: '0em',
    },
    h6: {
      fontSize: '1.1rem', // Used in training-form.tsx CardHeader
      fontWeight: 500,
      lineHeight: 1.4,
      letterSpacing: '0.0075em',
    },
    subtitle1: {
      fontSize: '1rem',
      fontWeight: 400,
      lineHeight: 1.5,
      letterSpacing: '0.00938em',
    },
    body1: {
      fontSize: '1rem',
      fontWeight: 400,
      lineHeight: 1.6, // Increased for better readability
      letterSpacing: '0.00938em',
    },
    body2: {
      fontSize: '0.875rem',
      fontWeight: 400,
      lineHeight: 1.5,
      letterSpacing: '0.01071em',
    },
    button: {
      textTransform: 'none', // No ALL CAPS for buttons
      fontWeight: 500, // Buttons should have a bit more emphasis
    },
  },
  shape: {
    borderRadius: 8, // Default border radius for components like Paper, Card, etc.
  },
  components: {
    MuiAppBar: {
      defaultProps: {
        elevation: 1, // Subtle shadow for AppBar
        color: 'primary', // Use primary color for AppBar background
      },
      styleOverrides: {
        root: {
          // AppBar specific styles if needed
        },
      },
    },
    MuiButton: {
      defaultProps: {
        // variant: 'contained', // Uncomment if most buttons should be contained
        // color: 'primary', // Uncomment if most buttons should use primary color
      },
      styleOverrides: {
        root: {
          borderRadius: 6, // Slightly less rounded than global shape.borderRadius for a more button-like feel
          padding: '8px 16px', // Default padding for buttons
        },
        containedSecondary: { // Specific style for secondary contained buttons
          color: '#000000', // Ensure contrast for yellow background
        }
      },
    },
    MuiCard: {
      defaultProps: {
        elevation: 2, // Default elevation for cards
      },
      styleOverrides: {
        root: {
          borderRadius: 12, // More rounded cards than default shape.borderRadius
        },
      },
    },
    MuiPaper: {
      defaultProps: {
        elevation: 1, // Default elevation for Paper components
      },
      styleOverrides: {
        root: {
          // General Paper styles
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        variant: 'outlined',
        size: 'small', // Make form inputs a bit more compact by default
      },
    },
    MuiSelect: {
      defaultProps: {
        variant: 'outlined',
        size: 'small',
      },
    },
    MuiFormControl: {
      defaultProps: {
        // variant: 'outlined', // This is handled by TextField/Select
      }
    },
    MuiDialog: {
      defaultProps: {
        PaperProps: {
          // elevation: 5, // Dialogs usually have higher elevation
        },
      },
      styleOverrides: {
        paper: { // Targeting the Paper component inside Dialog
          borderRadius: 10,
        },
      },
    },
    MuiLink: {
      defaultProps: {
        underline: 'hover',
      },
      styleOverrides: {
        root: {
          color: '#00796b', // Primary color for links
          fontWeight: 500,
          '&:hover': {
            color: '#004d40', // Darker shade on hover
          },
        },
      },
    },
    MuiCssBaseline: { // Optional: Define global scrollbar styles
      styleOverrides: `
        html {
          -webkit-font-smoothing: antialiased;
          -moz-osx-font-smoothing: grayscale;
          box-sizing: border-box;
        }
        body {
          overflow-x: hidden; /* Prevent horizontal scroll */
        }
        /* Basic scrollbar styling */
        ::-webkit-scrollbar {
          width: 8px;
          height: 8px;
        }
        ::-webkit-scrollbar-track {
          background: #f1f1f1;
        }
        ::-webkit-scrollbar-thumb {
          background: #c1c1c1;
          border-radius: 4px;
        }
        ::-webkit-scrollbar-thumb:hover {
          background: #a1a1a1;
        }
      `,
    },
  },
});

export default theme;
