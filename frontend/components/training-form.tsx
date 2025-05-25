"use client"

import { useState } from "react"
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardActions,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormGroup,
  FormControlLabel,
  Checkbox,
  CircularProgress,
} from "@mui/material"
import RecommendationsAPI from "@/features/recommendations/api"
import type { RecommendationRequest } from "@/features/recommendations/types"
import RecommendationDisplay from "@/components/recommendation-display"

const trainingGoals = [
  { id: "muscle-building", label: "筋肥大（ボディメイク）", value: "筋肥大" },
  { id: "fat-loss", label: "ダイエット（脂肪燃焼）", value: "ダイエット" },
  { id: "health", label: "健康維持", value: "健康維持" },
  { id: "performance", label: "パフォーマンス向上（スポーツ向け）", value: "パフォーマンス向上" },
]

const bodyParts = [
  { id: "full-body", label: "全身", value: "全身" },
  { id: "chest", label: "胸（大胸筋）", value: "胸" },
  { id: "back", label: "背中（広背筋・僧帽筋）", value: "背中" },
  { id: "shoulders", label: "肩（三角筋）", value: "肩" },
  { id: "arms", label: "腕（上腕二頭筋・上腕三頭筋）", value: "腕" },
  { id: "legs", label: "脚（大腿四頭筋・ハムストリング・ふくらはぎ）", value: "脚" },
  { id: "abs", label: "腹筋", value: "腹筋" },
]

const experienceLevels = [
  { id: "beginner", label: "初心者（1年未満）", value: "初心者" },
  { id: "intermediate", label: "中級者（1〜3年）", value: "中級者" },
  { id: "advanced", label: "上級者（3年以上）", value: "上級者" },
]

export default function TrainingForm() {
  const [trainingGoal, setTrainingGoal] = useState<string>("")
  const [selectedParts, setSelectedParts] = useState<string[]>([])
  const [experienceLevel, setExperienceLevel] = useState<string>("")
  const [availableTime, setAvailableTime] = useState<number>(30)
  const [isLoading, setIsLoading] = useState(false)
  const [recommendations, setRecommendations] = useState<string>("")

  const handlePartChange = (value: string) => {
    setSelectedParts((prev) => {
      if (prev.includes(value)) {
        return prev.filter((item) => item !== value)
      } else {
        return [...prev, value]
      }
    })
  }

  const handleSubmit = async () => {
    setIsLoading(true)

    try {
      const request: RecommendationRequest = {
        training_goal: trainingGoal,
        target_parts: selectedParts,
        experience_level: experienceLevel,
        available_time: availableTime,
      }

      console.log(request)

      const response = await RecommendationsAPI.proposeTrainingMenu(request)
      console.log(response)
      setRecommendations(response.recommendation)
    } catch (error) {
      console.error("トレーニングメニューの取得に失敗しました", error)
      alert("トレーニングメニューの取得に失敗しました。もう一度お試しください。")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Box sx={{ maxWidth: 800, margin: "0 auto", p: 2 }}>
      <Card>
        <CardHeader
          title="トレーニングメニュー提案"
          subheader="あなたの目標と条件に合わせたトレーニングメニューを提案します"
        />
        <CardContent>
          {/* トレーニング目的（プルダウン） */}
          <Box sx={{ mb: 3 }}>
            <FormControl fullWidth>
              <InputLabel id="training-goal-label">トレーニング目的</InputLabel>
              <Select
                labelId="training-goal-label"
                value={trainingGoal}
                label="トレーニング目的"
                onChange={(e) => setTrainingGoal(e.target.value)}
              >
                {/* <MenuItem value="">
                  <em>選択しない</em>
                </MenuItem> */}
                {trainingGoals.map((goal) => (
                  <MenuItem key={goal.id} value={goal.value}>
                    {goal.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>

          {/* トレーニングしたい部位（チェックボックス） */}
          <Box sx={{ mb: 3 }}>
            <Typography variant="subtitle1" gutterBottom>
              トレーニングしたい部位（複数選択可）
            </Typography>
            <FormGroup>
              {bodyParts.map((part) => (
                <FormControlLabel
                  key={part.id}
                  control={
                    <Checkbox
                      checked={selectedParts.includes(part.value)}
                      onChange={() => handlePartChange(part.value)}
                    />
                  }
                  label={part.label}
                />
              ))}
            </FormGroup>
          </Box>

          {/* トレーニング経験（プルダウン） */}
          <Box sx={{ mb: 3 }}>
            <FormControl fullWidth>
              <InputLabel id="experience-level-label">トレーニング経験</InputLabel>
              <Select
                labelId="experience-level-label"
                value={experienceLevel}
                label="トレーニング経験"
                onChange={(e) => setExperienceLevel(e.target.value)}
              >
                {/* <MenuItem value="">
                  <em>選択しない</em>
                </MenuItem> */}
                {experienceLevels.map((level) => (
                  <MenuItem key={level.id} value={level.value}>
                    {level.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>

          {/* トレーニング時間（プルダウン） */}
          <Box sx={{ mb: 3 }}>
            <FormControl fullWidth>
              <InputLabel id="available-time-label">トレーニング時間</InputLabel>
              <Select
                labelId="available-time-label"
                value={availableTime.toString()}
                label="トレーニング時間"
                onChange={(e) => setAvailableTime(Number(e.target.value))}
              >
                {[30, 60, 90, 120, 150, 180].map((time) => (
                  <MenuItem key={time} value={time.toString()}>
                    {time}分
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
        </CardContent>
        <CardActions>
          <Button onClick={handleSubmit} disabled={isLoading} variant="contained" fullWidth>
            トレーニングメニューを提案
          </Button>
        </CardActions>
      </Card>

      {isLoading && (
        <Box sx={{ display: "flex", justifyContent: "center", mt: 2 }}>
          <CircularProgress />
        </Box>
      )}

      {!isLoading && recommendations.length > 0 && (
        <Box sx={{ mt: 2 }}>
          <RecommendationDisplay recommendation={recommendations} />
        </Box>
      )}
    </Box>
  )
}
