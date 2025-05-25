import React, { ReactNode } from 'react';
import Link from 'next/link';
import { AppBar, Toolbar, Typography, Button, Container } from '@mui/material';

interface LayoutProps {
  children: ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <>
      {/* AppBar will use theme's defaultProps: elevation: 1, color: 'primary' */}
      <AppBar position="static"> 
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            <Link 
              href="/" 
              passHref 
              style={{ textDecoration: 'none', color: 'inherit' }}
              // sx={{ 
              //   textDecoration: 'none', 
              //   color: 'inherit', // Inherits color from AppBar (primary.contrastText)
              //   '&:hover': {
              //     textDecoration: 'underline', // Optional: underline on hover
              //   }
              // }}
            >
              Fitness Tracker
            </Link>
          </Typography>
          {/* Buttons will use theme's defaultProps (e.g., textTransform: 'none') */}
          {/* color="inherit" ensures they use primary.contrastText from the AppBar */}
          <Link href="/" passHref>
            <Button color="inherit">Home</Button>
          </Link>
          <Link href="/workouts" passHref>
            <Button color="inherit">Workouts</Button>
          </Link>
          <Link href="/recommendations" passHref>
            <Button color="inherit">Recommendations</Button>
          </Link>
        </Toolbar>
      </AppBar>
      {/* Container styling is fine, provides consistent page padding */}
      <Container sx={{ mt: { xs: 2, sm: 3, md: 4 }, mb: { xs: 2, sm: 3, md: 4 }, p: { xs: 1, sm: 2} }}>
        {children}
      </Container>
    </>
  );
};

export default Layout;
