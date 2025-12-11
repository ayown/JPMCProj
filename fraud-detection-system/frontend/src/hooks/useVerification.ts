import { useDispatch, useSelector } from 'react-redux';
import { verifyMessage, fetchHistory, fetchStats, clearCurrentVerification } from '@/store/verificationSlice';
import { RootState, AppDispatch } from '@/store';
import { VerificationRequest } from '@/types/verification';

export const useVerification = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { currentVerification, history, stats, isLoading, error } = useSelector(
    (state: RootState) => state.verification
  );

  const verify = async (data: VerificationRequest) => {
    return await dispatch(verifyMessage(data));
  };

  const loadHistory = (limit?: number, offset?: number) => {
    dispatch(fetchHistory({ limit, offset }));
  };

  const loadStats = () => {
    dispatch(fetchStats());
  };

  const clearVerification = () => {
    dispatch(clearCurrentVerification());
  };

  return {
    currentVerification,
    history,
    stats,
    isLoading,
    error,
    verify,
    loadHistory,
    loadStats,
    clearVerification,
  };
};

