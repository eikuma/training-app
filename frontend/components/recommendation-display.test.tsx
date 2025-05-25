import React from 'react';
import { render, screen } from '@testing-library/react';
import RecommendationDisplay from './recommendation-display'; // Adjust path as necessary
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from '../app/theme'; // Import your actual theme

// No longer mocking ThemeProvider

describe('RecommendationDisplay Component', () => {
  const renderWithTheme = (component) => {
    return render(
      <ThemeProvider theme={theme}>
        <CssBaseline />
        {component}
      </ThemeProvider>
    );
  };

  const mockRecommendation = "This is a mock recommendation text.\nIt includes exercises like:\n- Push-ups\n- Squats\n- Plank";

  test('renders recommendation text correctly', () => {
    renderWithTheme(<RecommendationDisplay recommendation={mockRecommendation} />);
    
    expect(screen.getByText(/This is a mock recommendation text/i)).toBeInTheDocument();
    expect(screen.getByText(/- Push-ups/i)).toBeInTheDocument();
    expect(screen.getByText(/- Squats/i)).toBeInTheDocument();
    expect(screen.getByText(/- Plank/i)).toBeInTheDocument();
  });

  test('renders card header title', () => {
    renderWithTheme(<RecommendationDisplay recommendation={mockRecommendation} />);

    // The title is within a span styled as h6. Query by text content.
    const headerTitle = screen.getByText(/あなたへのおすすめトレーニングメニュー/i);
    expect(headerTitle).toBeInTheDocument();
    // Optionally, check its class if styles are important for the test
    expect(headerTitle).toHaveClass('MuiTypography-h6'); // From titleTypographyProps={{ variant: 'h6' }}
  });

  test('renders nothing if recommendation prop is null or empty', () => {
    const { container: containerNull } = renderWithTheme(<RecommendationDisplay recommendation={null} />);
    expect(containerNull).toBeEmptyDOMElement();

    const { container: containerEmpty } = renderWithTheme(<RecommendationDisplay recommendation="" />);
    expect(containerEmpty).toBeEmptyDOMElement();
  });
});
