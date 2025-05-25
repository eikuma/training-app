'use client';
import React from "react";
import {
  Card,
  CardHeader,
  CardContent,
  Typography
} from "@mui/material";
import { useTheme } from '@mui/material/styles';

import { RecommendationResponse } from '@/features/recommendations/types/index';

export default function RecommendationDisplay({ recommendation }: RecommendationResponse) {
  const theme = useTheme();
  if (!recommendation) return null;

  return (
    // Card will use theme defaults (elevation=2, borderRadius=12) unless overridden by sx here.
    <Card 
      elevation={3} // Explicitly setting higher elevation for emphasis as a result display
      sx={{ 
        width: "100%", 
        mt: 3, // Keep increased margin top from previous refactor attempt
        // Use a slightly off-white for paper in light mode, or a dark grey for dark mode.
        // theme.palette.background.paper is usually #fff, this provides a subtle difference.
        backgroundColor: theme.palette.mode === 'dark' ? theme.palette.grey[800] : theme.palette.grey[50], 
        borderColor: 'primary.main', // Use theme's primary color string directly
        borderTopWidth: 3,
        borderTopStyle: 'solid',
        // borderRadius will be from theme's MuiCard default (12px)
      }}
    >
      <CardHeader
        title="あなたへのおすすめトレーニングメニュー" // Use the more personalized title
        titleTypographyProps={{ 
          variant: 'h6', // Theme h6 will be applied
          align: 'center', 
          color: 'primary.main' // Emphasize title with primary color
        }}
        subheader="以下のメニューを試して、目標達成を目指しましょう！" // Use the more engaging subheader
        subheaderTypographyProps={{ 
          variant: 'subtitle1', // Theme subtitle1 will be applied
          align: 'center', 
          color: 'text.secondary', // Use theme's secondary text color
          sx: { mb: 0 } // Remove bottom margin for subheader if any from variant
        }}
        sx={{pb: 1}} // Adjust padding for CardHeader
      />
      <CardContent sx={{pt: 1}}> {/* Adjust padding for CardContent */}
        {/* Typography will use theme's body1. whiteSpace and lineHeight are specific formatting. */}
        <Typography 
          variant="body1" // Theme body1 will be applied
          sx={{ 
            whiteSpace: 'pre-wrap', 
            textAlign: 'left', 
            p:1, 
            lineHeight: 1.7, // Keep improved line height
            color: 'text.primary' // Use theme's primary text color
          }}
        >
          {recommendation}
        </Typography>
      </CardContent>
    </Card>
  );
}
