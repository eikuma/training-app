export type WorkoutSession = {
    id:     number;
	date:   string;
	user_id: number;
};

export type Exercise = {
    exercise_id:           number;
    session_id:    number; 
    exercise_name :string;
    sets?:        Set[];
};

export type Set = {
    set_id:         number;
    exercise_id: number;
    set_number:  number;
    weight:     number;
    reps:       number;
};

export type WorkoutsResponse = {
    workouts: WorkoutSession[];
}

export type Workout = {
    id:     number;
	date:   string;
	userId: string;
    exercises?: Exercise[];
}

export type WorkoutResponse = {
    workout: Workout;
}


export type WorkoutsParams = {
    id?: number;
    date?: string;
}

export type WorkoutSessionRequest = {
    date: string;
    user_id: number;
}

export type ExerciseRequest = {
    exercise_name: string;
}

export type ExerciseResponse = {
    exercise: Exercise;
}

export type SetRequest = {
    set_number: number;
    weight: number;
    reps: number;
}

export type SetResponse = {
    sets: Set[];
}
