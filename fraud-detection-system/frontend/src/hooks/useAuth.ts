import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { login, register, logout, getProfile } from '@/store/authSlice';
import { RootState, AppDispatch } from '@/store';
import { LoginRequest, RegisterRequest } from '@/types/auth';

export const useAuth = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth);

  const handleLogin = async (data: LoginRequest) => {
    const result = await dispatch(login(data));
    if (login.fulfilled.match(result)) {
      navigate('/dashboard');
    }
  };

  const handleRegister = async (data: RegisterRequest) => {
    const result = await dispatch(register(data));
    if (register.fulfilled.match(result)) {
      navigate('/login');
    }
  };

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  const refreshProfile = () => {
    dispatch(getProfile());
  };

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    login: handleLogin,
    register: handleRegister,
    logout: handleLogout,
    refreshProfile,
  };
};

