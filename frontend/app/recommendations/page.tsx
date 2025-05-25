import React from 'react';
import TrainingForm from "@/components/training-form";
import { Container, Typography, Box } from '@mui/material';

const RecommendationsPage = () => {
  return (
    // mt, mb removed as Layout.tsx now provides responsive vertical spacing for its Container.
    // Page-specific mb on the Box below is kept.
    <Container maxWidth="md"> 
      <Box sx={{ textAlign: 'center', mb: { xs: 3, sm: 4 } }}> {/* Responsive margin bottom */}
        <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'text.primary' }}>
          パーソナライズドトレーニング提案
        </Typography>
        <Typography variant="subtitle1" color="text.secondary"> {/* Changed to text.secondary */}
          あなたの目標や経験に合わせて、最適なトレーニングメニューをご提案します。
        </Typography>
      </Box>
      <TrainingForm />
    </Container>
  );
};

export default RecommendationsPage;

