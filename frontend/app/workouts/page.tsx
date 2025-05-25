"use client";

import React, { useState } from "react";
import Link from 'next/link';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardActions,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Typography,
  TextField,
  CircularProgress,
} from "@mui/material";
import WorkoutsAPI from "@/features/workouts/api";
import type {
  WorkoutsParams,
  WorkoutsResponse,
  WorkoutResponse,
  WorkoutSessionRequest,
  ExerciseRequest,
  ExerciseResponse,
  SetRequest,
  SetResponse,
  Exercise,
  Set,
} from "@/features/workouts/types";

// 例として利用するトレーニングメニュー一覧
const trainingMenus = [
  "ベンチプレス",
  "スクワット",
  "デッドリフト",
  "ショルダープレス",
  "ラットプルダウン",
];

// 重量の選択肢：0.5kg刻み（0.5～100kg）
const weightOptions: number[] = [];
for (let w = 0.5; w <= 100; w += 0.5) {
  weightOptions.push(parseFloat(w.toFixed(1)));
}

const Workouts = () => {
  // 日付入力（yyyy-mm-dd形式）
  const [date, setDate] = useState<string>("");
  // セッション情報（取得済みの場合）
  const [session, setSession] = useState<WorkoutResponse | null>(null);
  // セッション内のトレーニングメニュー（Exercise）一覧
  const [exercises, setExercises] = useState<Exercise[]>([]);
  // ローディング状態
  const [loading, setLoading] = useState<boolean>(false);
  // 「トレーニングメニュー追加」フォームの表示
  const [showExerciseForm, setShowExerciseForm] = useState<boolean>(false);
  // ドロップダウンで選択中のトレーニングメニュー
  const [selectedMenu, setSelectedMenu] = useState<string>("");
  // セット入力ダイアログの制御
  const [setDialogOpen, setSetDialogOpen] = useState<boolean>(false);
  // 現在、セットを追加する対象の exerciseId
  const [currentExerciseId, setCurrentExerciseId] = useState<number | null>(null);
  // セット入力用状態
  const [setNumber, setSetNumber] = useState<number>(1);
  const [weight, setWeight] = useState<number>(0.5);
  const [reps, setReps] = useState<number>(1);

  // 「トレーニングを始める」ボタン押下時の処理
  const handleStartTraining = async () => {
    if (!date) {
      alert("日付を入力してください");
      return;
    }
    setLoading(true);
    try {
      const isoDate = new Date(date).toISOString();
      const params: WorkoutsParams = { date: isoDate };
      const response: WorkoutsResponse = await WorkoutsAPI.fetchWorkouts(params);
      console.log("response", response);
      let sessionData: WorkoutResponse | null = null;
      if (response.workouts.length === 0) {
        const request: WorkoutSessionRequest = {
          date: isoDate,
          user_id: -1, // ユーザーIDは実際の値に置換
        };
        sessionData = await WorkoutsAPI.createWorkoutSession(request);
      } else {
        const sessionId = response.workouts[0].id;
        sessionData = await WorkoutsAPI.fetchWorkout(sessionId);
      }
      setExercises(sessionData.workout?.exercises ?? []);
      setSession(sessionData);
      setShowExerciseForm(true);
    } catch (error) {
      console.error("セッションの取得または作成に失敗しました", error);
      alert("セッションの取得または作成に失敗しました。");
    } finally {
      setLoading(false);
    }
  };

  // トレーニングメニュー追加フォームの送信
  const handleAddExercise = async () => {
    if (!selectedMenu || !session) {
      alert("トレーニングメニューを選択してください");
      return;
    }
    setLoading(true);
    try {
      const request: ExerciseRequest = {
        exercise_name: selectedMenu,
      };
      const newExercise: ExerciseResponse = await WorkoutsAPI.createExercise(
        session.workout.id,
        request
      );
      console.log("newExercise:", newExercise);
      setExercises((prev) => [newExercise.exercise, ...prev]);
      setSelectedMenu("");
    } catch (error) {
      console.error("トレーニングメニューの追加に失敗しました", error);
      alert("トレーニングメニューの追加に失敗しました。");
    } finally {
      setLoading(false);
    }
  };

  // セット入力ダイアログを開く（対象の exerciseId を保持）
  const handleOpenSetDialog = (exerciseId: number) => {
    setCurrentExerciseId(exerciseId);
    setSetNumber(1);
    setWeight(0.5);
    setReps(1);
    setSetDialogOpen(true);
  };

  // セット入力ダイアログの送信処理
  const handleSubmitSet = async () => {
    if (!session || currentExerciseId === null) {
      return;
    }
    try {
      const request: SetRequest = {
        set_number: setNumber,
        weight,
        reps,
      };
      const response: SetResponse = await WorkoutsAPI.createSet(
        session.workout.id,
        currentExerciseId,
        request
      );
      setExercises((prevExercises) =>
        prevExercises.map((ex) =>
          ex.exercise_id === currentExerciseId ? { ...ex, sets: response.sets } : ex
        )
      );
      setSetDialogOpen(false);
    } catch (error) {
      console.error("セットの追加に失敗しました", error);
      alert("セットの追加に失敗しました。");
    }
  };

  return (
    // Box p:2, mt:2, mb:2 are specific to this page, can be kept. Layout.tsx handles broader page padding.
    <Box sx={{ maxWidth: 800, mx: "auto", p: {xs: 1, sm: 2}, mt: {xs: 1, sm: 2}, mb: {xs:1, sm:2} }}>
      {/* 日付入力と「トレーニングを始める」ボタン */}
      {!session && (
        // elevation={3} is fine, overrides default Paper elevation (1)
        <Paper elevation={3} sx={{ mb: 3, p: { xs: 2, sm: 3 } }}> 
          <Typography variant="h5" component="h2" sx={{ mb: 2, textAlign: 'center', color: 'text.primary' }}>
            トレーニング記録を開始
          </Typography>
          {/* TextField will use theme's defaultProps: variant="outlined", size="small" */}
          <TextField
            label="日付を選択"
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            fullWidth
            InputLabelProps={{ shrink: true }}
            sx={{ mb: 2 }}
          />
          {/* Button will use theme's defaultProps: textTransform, borderRadius */}
          <Button
            variant="contained"
            color="primary" // Uses theme.palette.primary.main
            onClick={handleStartTraining}
            disabled={loading || !date}
            fullWidth
            size="large" // Specific size choice
            startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}
          >
            {loading ? "準備中..." : "トレーニングを始める"}
          </Button>
        </Paper>
      )}

      {/* セッション開始後の表示 */}
      {session && (
        <Box>
          {/* elevation={2} is fine, overrides default Paper elevation (1) */}
          <Paper elevation={2} sx={{ mb: 3, p: { xs: 1.5, sm: 2.5 } }}> 
            <Typography variant="h5" component="h2" gutterBottom sx={{textAlign: 'center', color: 'text.primary'}}>
              {session.workout.date 
                ? new Date(session.workout.date).toLocaleDateString('ja-JP', { year: 'numeric', month: 'long', day: 'numeric' })
                : 'トレーニングセッション'}
            </Typography>
          </Paper>

          {/* トレーニングメニュー追加フォーム */}
          {showExerciseForm && (
            // elevation={2} is fine
            <Paper elevation={2} sx={{ mb: 3, p: { xs: 2, sm: 3 } }}> 
              <Typography variant="h6" component="h3" sx={{ mb: 2, color: 'text.primary' }}>
                エクササイズを追加
              </Typography>
              {/* FormControl and Select will use theme's defaultProps */}
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel id="training-menu-label">エクササイズを選択</InputLabel>
                <Select
                  labelId="training-menu-label"
                  value={selectedMenu}
                  label="エクササイズを選択" // Updated label to match InputLabel
                  onChange={(e) => setSelectedMenu(e.target.value)}
                >
                  {trainingMenus.map((menu) => (
                    <MenuItem key={menu} value={menu}>
                      {menu}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Button
                variant="contained"
                color="secondary" // Uses theme.palette.secondary.main
                onClick={handleAddExercise}
                disabled={loading || !selectedMenu}
                startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}
              >
                {loading ? "追加中..." : "エクササイズを追加"}
              </Button>
            </Paper>
          )}

          {/* トレーニングメニュー一覧表示 */}
          {exercises.length > 0 ? (
            exercises.map((exercise) => (
              // Card will use theme's defaultProps: elevation=2, borderRadius=12
              <Card key={exercise.exercise_id} sx={{ mb: 2, p: 1 }}> 
                <CardHeader
                  title={exercise.exercise_name}
                  titleTypographyProps={{ variant: 'h6' }} // Theme h6 styles applied
                  sx={{ pb: 0, '& .MuiCardHeader-title': { color: 'text.primary' } }} // Ensure title color
                />
                <CardContent>
                  {exercise.sets && exercise.sets.length > 0 ? (
                    <Box sx={{ display: "flex", overflowX: "auto", gap: 1.5, py: 1 }}>
                      {exercise.sets.map((s: Set, index: number) => (
                        // elevation={2} is fine for these inner Paper elements, makes them pop a bit
                        <Paper 
                          key={s.set_id} 
                          elevation={2} 
                          sx={{ 
                            minWidth: { xs: 110, sm: 130 }, 
                            p: 1.5, 
                            textAlign: 'center',
                            // Using theme's grey shade for background
                            backgroundColor: (theme) => theme.palette.mode === 'dark' ? theme.palette.grey[700] : theme.palette.grey[100],
                          }}
                        >
                          <Typography variant="subtitle2" gutterBottom sx={{ color: 'text.primary' }}>
                            {index + 1} セット
                          </Typography>
                          <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                            {s.weight} kg
                          </Typography>
                          <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                            {s.reps} 回
                          </Typography>
                        </Paper>
                      ))}
                    </Box>
                  ) : (
                    <Typography variant="body2" color="text.secondary" sx={{mt: 1, mb: 1, textAlign: 'center'}}>
                      このエクササイズのセットはまだ記録されていません。
                    </Typography>
                  )}
                </CardContent>
                <CardActions sx={{ justifyContent: 'flex-end', pt:0, pr: 1.5, pb: 1.5 }}>
                  <Button
                    variant="contained"
                    size="small" // Specific size choice
                    color="primary"
                    onClick={() => handleOpenSetDialog(exercise.exercise_id)}
                  >
                    ＋ セット追加
                  </Button>
                </CardActions>
              </Card>
            ))
          ) : (
            showExerciseForm && (
              <Typography variant="body1" color="text.secondary" sx={{ textAlign: 'center', mt: 3 }}>
                上記フォームからエクササイズを追加してください。
              </Typography>
            )
          )}
        </Box>
      )}

      {/* Dialog will use theme's defaultProps for Paper (borderRadius) */}
      <Dialog open={setDialogOpen} onClose={() => { setSetDialogOpen(false); setCurrentExerciseId(null);}} fullWidth maxWidth="xs">
        <DialogTitle sx={{ textAlign: 'center', color: 'text.primary' }}>セット情報を追加</DialogTitle>
        {/* pt: '16px !important' is a bit of a hack, prefer theme adjustments if possible, but can keep if necessary */}
        <DialogContent sx={{ pt: '16px !important' }}> 
          <Box component="form" sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <FormControl fullWidth>
              <InputLabel id="set-number-label">セット番号</InputLabel>
              {/* Select will use theme's defaultProps */}
            <Select
              labelId="set-number-label"
              value={setNumber.toString()}
              label="セット番号"
              onChange={(e) => setSetNumber(Number(e.target.value))}
            >
              {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((n) => (
                <MenuItem key={n} value={n}>
                  {n}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl fullWidth>
            <InputLabel id="weight-label">重量 (kg)</InputLabel>
            <Select
              labelId="weight-label"
              value={weight.toString()}
              label="重量 (kg)"
              onChange={(e) => setWeight(Number(e.target.value))}
            >
              {weightOptions.map((w) => (
                <MenuItem key={w} value={w}>
                  {w} kg
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl fullWidth>
            <InputLabel id="reps-label">回数</InputLabel>
            <Select
              labelId="reps-label"
              value={reps.toString()}
              label="回数"
              onChange={(e) => setReps(Number(e.target.value))}
            >
              {Array.from({ length: 30 }, (_, i) => i + 1).map((r) => (
                <MenuItem key={r} value={r}>
                  {r} 回
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          </Box>
        </DialogContent>
        <DialogActions sx={{ p: '16px 24px' }}> {/* Padding is fine */}
          {/* Cancel button text color will be default (primary by Mui default, or can be styled) */}
          <Button onClick={() => { setSetDialogOpen(false); setCurrentExerciseId(null);}}>キャンセル</Button>
          <Button variant="contained" onClick={handleSubmitSet} disabled={loading}>
            {loading ? <CircularProgress size={20} color="inherit" /> : "記録する"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Workouts;
