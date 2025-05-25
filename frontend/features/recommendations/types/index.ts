export type Recommendation = {
    training_goal: string;
    target_parts: string[];
    experience_level: string;
    available_time: number;
};

export type RecommendationRequest = Pick<Recommendation, 'training_goal' | 'target_parts' | 'experience_level' | 'available_time'>;

export type RecommendationResponse = {
    recommendation: string;
};
