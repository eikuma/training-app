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
    const token = localStorage.getItem('authToken');
    const config: any = { params }; // Start with params for GET request

    if (token) {
      config.headers = {
        ...config.headers, // Spread any existing headers from apiClient default config if needed
        'Authorization': `Bearer ${token}`,
      };
    }

    try {
      console.log('Fetching workouts with params:', params, 'and config:', config);
      const { data } = await apiClient.get<WorkoutsResponse>('/workouts', config);
      return data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        console.error("Unauthorized: Token might be invalid or expired. (fetchWorkouts)");
        // Optionally redirect to login or clear token
        // localStorage.removeItem('authToken');
        // window.location.href = '/login';
      }
      throw error; // Re-throw to be handled by the calling component
    }
  },
  async fetchWorkout(id: number): Promise<WorkoutResponse> {
    // Assuming fetching a single workout might also require auth, add token if needed
    const token = localStorage.getItem('authToken');
    const config: any = {};
    if (token) {
      config.headers = { 'Authorization': `Bearer ${token}` };
    }
    try {
      const { data } = await apiClient.get<WorkoutResponse>(`/workouts/${id}`, config);
      return data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        console.error("Unauthorized: Token might be invalid or expired. (fetchWorkout)");
      }
      throw error;
    }
  },
  async createWorkoutSession(body: WorkoutSessionRequest): Promise<WorkoutResponse> {
    const token = localStorage.getItem('authToken');
    const config: any = {};

    if (token) {
      config.headers = {
        // Content-Type is already set in apiClient default config
        'Authorization': `Bearer ${token}`,
      };
    }
    // The user_id in the 'body' will be ignored by the backend if present,
    // as it uses the user_id from the token.

    try {
      const { data } = await apiClient.post<WorkoutResponse>('/workouts', body, config);
      return data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        console.error("Unauthorized: Token might be invalid or expired. (createWorkoutSession)");
        // Optionally redirect to login or clear token
        // localStorage.removeItem('authToken');
        // window.location.href = '/login';
      }
      throw error; // Re-throw to be handled by the calling component
    }
  },
  async createExercise(sessionId: number, body: ExerciseRequest): Promise<ExerciseResponse> {
    // Assuming creating an exercise also requires auth
    const token = localStorage.getItem('authToken');
    const config: any = {};
    if (token) {
      config.headers = { 'Authorization': `Bearer ${token}` };
    }
    try {
      const { data } = await apiClient.post<ExerciseResponse>(`/workouts/${sessionId}/exercises`, body, config);
      return data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        console.error("Unauthorized: Token might be invalid or expired. (createExercise)");
      }
      throw error;
    }
  },
  async createSet(sessionId: number, exerciseId: number, body: SetRequest): Promise<SetResponse> {
    // Assuming creating a set also requires auth
    const token = localStorage.getItem('authToken');
    const config: any = {};
    if (token) {
      config.headers = { 'Authorization': `Bearer ${token}` };
    }
    try {
      const { data } = await apiClient.post<SetResponse>(`/workouts/${sessionId}/exercises/${exerciseId}/sets`, body, config);
      return data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        console.error("Unauthorized: Token might be invalid or expired. (createSet)");
      }
      throw error;
    }
  },
};

export default WorkoutsAPI;
