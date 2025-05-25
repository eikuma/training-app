import React from "react";
import {
  Card,
  CardHeader,
  CardContent,
  Typography
} from "@mui/material";

import { RecommendationResponse } from '@/features/recommendations/types/index';

export default function RecommendationDisplay({ recommendation }: RecommendationResponse) {
  if (!recommendation) return null;

  return (
    <Card sx={{ width: "100%", mt: 2 }}>
      <CardHeader
        title="おすすめトレーニングメニュー"
        subheader="あなたの目標と条件に合わせたトレーニングメニューです"
      />
      <CardContent>
        <Typography sx={{ whiteSpace: 'pre-wrap' }}>
          {recommendation}
        </Typography>
      </CardContent>
    </Card>
  );
}
