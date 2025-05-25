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
    <Box sx={{ maxWidth: 800, mx: "auto", p: 2 }}>
      <div style={{ margin: "1rem" }}>
        <Link href="/" passHref>
          <button style={{ padding: "0.5rem 1rem", fontSize: "1rem" }}>ホームへ</button>
        </Link>
      </div>
      {/* 日付入力と「トレーニングを始める」ボタン */}
      {!session && (
        <Card sx={{ mb: 3, p: 2 }}>
          <Typography variant="h6" sx={{ mb: 2 }}>
            トレーニング記録を開始
          </Typography>
          <TextField
            label="日付"
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            fullWidth
            InputLabelProps={{ shrink: true }}
            sx={{ mb: 2 }}
          />
          <Button
            variant="contained"
            onClick={handleStartTraining}
            disabled={loading}
            fullWidth
          >
            {loading ? <CircularProgress size={24} /> : "トレーニングを始める"}
          </Button>
        </Card>
      )}

      {/* セッション開始後の表示 */}
      {session && (
        <Box>
          <Card sx={{ mb: 2, p: 2 }}>
            <CardHeader
              title={`日付: ${session.workout.date}`}
            />
            <CardContent>
              <Typography variant="body1">
                トレーニングメニューが以下に表示されます。
              </Typography>
            </CardContent>
          </Card>

          {/* トレーニングメニュー追加フォーム */}
          {showExerciseForm && (
            <Card sx={{ mb: 3, p: 2 }}>
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel id="training-menu-label">トレーニングメニュー</InputLabel>
                <Select
                  labelId="training-menu-label"
                  value={selectedMenu}
                  label="トレーニングメニュー"
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
                onClick={handleAddExercise}
                disabled={loading}
              >
                追加する
              </Button>
            </Card>
          )}

          {/* トレーニングメニュー一覧表示 */}
          {exercises.map((exercise) => (
            <Card key={exercise.exercise_id} sx={{ mb: 2, p: 2 }}>
              <CardHeader
                title={exercise.exercise_name}
              />
              <CardContent>
                {/* セット一覧は横スクロール */}
                <Box sx={{ display: "flex", overflowX: "auto", gap: 1, pb: 1 }}>
                  {exercise.sets &&
                    exercise.sets.map((s: Set, index: number) => (
                      <Card key={s.set_id} sx={{ minWidth: 120, p: 1 }}>
                        <Typography variant="caption">
                          {index + 1}セット
                        </Typography>
                        <Typography variant="body2">
                          重量: {s.weight}kg
                        </Typography>
                        <Typography variant="body2">
                          回数: {s.reps}
                        </Typography>
                      </Card>
                    ))}
                </Box>
              </CardContent>
              <CardActions>
                <Button
                  variant="outlined"
                  onClick={() => handleOpenSetDialog(exercise.exercise_id)}
                >
                  ＋ セット追加
                </Button>
              </CardActions>
            </Card>
          ))}
        </Box>
      )}

      {/* セット入力ダイアログ */}
      <Dialog open={setDialogOpen} onClose={() => setSetDialogOpen(false)}>
        <DialogTitle>セット情報を追加</DialogTitle>
        <DialogContent
          sx={{ display: "flex", flexDirection: "column", gap: 2, mt: 1 }}
        >
          {/* セット番号 */}
          <FormControl fullWidth>
            <InputLabel id="set-number-label">セット番号</InputLabel>
            <Select
              labelId="set-number-label"
              value={setNumber.toString()}
              label="セット番号"
              onChange={(e) => setSetNumber(Number(e.target.value))}
            >
              {[1, 2, 3, 4, 5].map((n) => (
                <MenuItem key={n} value={n}>
                  {n}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          {/* 重量 (kg) */}
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
                  {w}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          {/* 回数 */}
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
                  {r}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setSetDialogOpen(false)}>キャンセル</Button>
          <Button variant="contained" onClick={handleSubmitSet}>
            記録する
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Workouts;
