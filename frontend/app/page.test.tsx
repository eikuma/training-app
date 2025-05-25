import React from 'react';
import { render, screen } from '@testing-library/react';
import HomePage from './page'; // Adjust path to your actual home page component
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme'; // Import your actual theme

// No longer mocking ThemeProvider

describe('HomePage Component', () => {
  const renderWithTheme = (component) => {
    return render(
      <ThemeProvider theme={theme}>
        <CssBaseline />
        {component}
      </ThemeProvider>
    );
  };

  test('renders hero section elements', () => {
    renderWithTheme(<HomePage />);

    const welcomeHeading = screen.getByRole('heading', { name: /フィットネストラッカーへようこそ！/i });
    expect(welcomeHeading).toBeInTheDocument();

    const descriptionText = /あなたのトレーニングを記録し、パーソナライズされたおすすめメニューで目標達成をサポートします。/i;
    expect(screen.getByText(descriptionText)).toBeInTheDocument();
  });

  test('renders Call to Action (CTA) links with correct href attributes', () => {
    renderWithTheme(<HomePage />);

    // The buttons are wrapped in Link components, which are mocked as <a> tags.
    // We should query for links by their text content.
    const recordWorkoutLink = screen.getByRole('link', { name: /トレーニングを記録/i });
    expect(recordWorkoutLink).toBeInTheDocument();
    expect(recordWorkoutLink).toHaveAttribute('href', '/workouts');

    const getRecommendationLink = screen.getByRole('link', { name: /おすすめメニューを取得/i });
    expect(getRecommendationLink).toBeInTheDocument();
    expect(getRecommendationLink).toHaveAttribute('href', '/recommendations');
  });
});
