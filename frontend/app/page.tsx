'use client';
import React from 'react';
import Link from 'next/link';
import { Box, Typography, Button, Container, Paper } from '@mui/material';
import { useTheme } from '@mui/material/styles';

const Home: React.FC = () => {
  const theme = useTheme();
  return (
    // Container's mt/mb might be slightly adjusted by Layout.tsx's new responsive padding. This is fine.
    <Container maxWidth="md" sx={{ textAlign: 'center' /* Removed mt, mb as Layout.tsx handles it */ }}>
      <Paper
        elevation={3} // Explicitly keeping higher elevation for hero section
        sx={{
          padding: { xs: 3, sm: 4, md: 6 }, // Adjusted responsive padding slightly
          mb: 4, // Margin bottom for spacing before CTAs
          // Using theme's background.default for light mode, and a dark grey for dark mode
          backgroundColor: theme.palette.mode === 'dark' ? theme.palette.grey[800] : theme.palette.background.default,
        }}
      >
        <Typography
          variant="h3" // Theme provides base h3 styling (size, weight)
          component="h1"
          gutterBottom
          sx={{
            // fontWeight: 'bold', // Theme h3 has fontWeight: 600. If 'bold' (700) is truly needed, keep this. For now, try theme's weight.
            fontSize: { xs: '2rem', sm: '2.5rem', md: '3rem' }, // Keep responsive font size override
            color: 'text.primary', // Ensure it uses theme's primary text color
          }}
        >
          フィットネストラッカーへようこそ！
        </Typography>
        <Typography
          variant="body1" // Theme provides base body1 styling
          color="text.secondary" // Using theme's secondary text color
          sx={{
            mb: 3,
            fontSize: { xs: '0.9rem', sm: '1rem' }, // Keep responsive font size override
          }}
        >
          あなたのトレーニングを記録し、パーソナライズされたおすすめメニューで目標達成をサポートします。
        </Typography>
      </Paper>

      <Box
        sx={{
          display: 'flex',
          flexDirection: { xs: 'column', sm: 'row' },
          justifyContent: 'center',
          alignItems: 'center',
          gap: 2,
        }}
      >
        {/* Buttons will pick up default styles from theme (textTransform, borderRadius) */}
        <Link href="/workouts" passHref>
          <Button
            variant="contained" // This is standard, not from defaultProps unless set globally
            color="primary"     // Will use theme.palette.primary.main
            size="large"        // Specific choice
            sx={{ minWidth: { xs: '100%', sm: '200px' } }} // Responsive minWidth
          >
            トレーニングを記録
          </Button>
        </Link>
        <Link href="/recommendations" passHref>
          <Button
            variant="outlined" // This is standard
            color="primary"    // Will use theme.palette.primary.main
            size="large"       // Specific choice
            sx={{ minWidth: { xs: '100%', sm: '200px' } }} // Responsive minWidth
          >
            おすすめメニューを取得
          </Button>
        </Link>
      </Box>
    </Container>
  );
};

export default Home;
