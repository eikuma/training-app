import React from 'react';
import { render, screen } from '@testing-library/react';
import Layout from './Layout'; // Adjust path as necessary
import '@testing-library/jest-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from '../app/theme'; // Import your actual theme

// No longer mocking ThemeProvider and CssBaseline, using actual ones.

describe('Layout Component', () => {
  const renderWithTheme = (component) => {
    return render(
      <ThemeProvider theme={theme}>
        <CssBaseline />
        {component}
      </ThemeProvider>
    );
  };

  test('renders children correctly', () => {
    const childText = 'Test Child Content';
    renderWithTheme(
      <Layout>
        <div>{childText}</div>
      </Layout>
    );
    expect(screen.getByText(childText)).toBeInTheDocument();
  });

  test('renders navigation links with correct text and href attributes', () => {
    renderWithTheme(<Layout><div></div></Layout>);

    // Check for the title link
    const titleLink = screen.getByText('Fitness Tracker'); // This is the text within the <a> tag
    expect(titleLink).toBeInTheDocument();
    // The mock now renders an <a> tag directly if children is a string
    expect(titleLink).toHaveAttribute('href', '/'); 


    // Check for other navigation links, which are MUI Buttons wrapped in our Link mock
    // The mock clones props onto the Button, and the Button itself renders a button element.
    // However, the link behavior is what we test via the href on the surrounding <a>.
    // The text "Home" is inside a Button, which is inside an <a>.
    const homeLink = screen.getByRole('link', { name: /home/i });
    expect(homeLink).toBeInTheDocument();
    expect(homeLink).toHaveAttribute('href', '/');
    
    const workoutsLink = screen.getByRole('link', { name: /workouts/i });
    expect(workoutsLink).toBeInTheDocument();
    expect(workoutsLink).toHaveAttribute('href', '/workouts');

    const recommendationsLink = screen.getByRole('link', { name: /recommendations/i });
    expect(recommendationsLink).toBeInTheDocument();
    expect(recommendationsLink).toHaveAttribute('href', '/recommendations');
  });
});
