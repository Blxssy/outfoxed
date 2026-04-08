export type LoginRequest = {
    nickName: string;
    password: string;
};

export type RegisterRequest = {
    nickName: string;
    email: string;
    password: string;
    confirmPassword: string;
};
