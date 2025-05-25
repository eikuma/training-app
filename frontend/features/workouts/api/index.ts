import axios from 'axios';
import { 
    WorkoutsParams, 
    WorkoutsResponse, 
    WorkoutResponse,
    WorkoutSessionRequest,
    ExerciseRequest,
    ExerciseResponse,
    SetRequest,
    SetResponse
} from '@/features/workouts/types';

const API_URL = "http://localhost:8080";
const apiClient = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});


interface IWorkoutsAPI {
    fetchWorkouts(params:WorkoutsParams): Promise<WorkoutsResponse>;
    fetchWorkout(id: number):Promise<WorkoutResponse>;
    createWorkoutSession(body:WorkoutSessionRequest):Promise<WorkoutResponse>;
    createExercise(sessionId:number, body:ExerciseRequest):Promise<ExerciseResponse>;
    createSet(sessionId: number,exerciseId:number, body:SetRequest):Promise<SetResponse>;
}

const WorkoutsAPI: IWorkoutsAPI = {
  async fetchWorkouts(params: WorkoutsParams): Promise<WorkoutsResponse> {
    console.log(params);
    const { data } = await apiClient.get<WorkoutsResponse>('/workouts', { params });
    return data;
  },
  async fetchWorkout(id: number): Promise<WorkoutResponse> {
    const { data } = await apiClient.get<WorkoutResponse>(`/workouts/${id}`);
    return data;
  },
  async createWorkoutSession(body: WorkoutSessionRequest): Promise<WorkoutResponse> {
    const { data } = await apiClient.post<WorkoutResponse>('/workouts', body);
    return data;
  },
  async createExercise(sessionId: number, body: ExerciseRequest): Promise<ExerciseResponse> {
    const { data } = await apiClient.post<ExerciseResponse>(`/workouts/${sessionId}/exercises`, body);
    return data;
  },
  async createSet(sessionId: number, exerciseId: number, body: SetRequest): Promise<SetResponse> {
    const { data } = await apiClient.post<SetResponse>(`/workouts/${sessionId}/exercises/${exerciseId}/sets`, body);
    return data;
  },
};

export default WorkoutsAPI;
