import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { verificationService } from '@/services/verification';
import { VerificationRequest, VerificationResponse, VerificationStats } from '@/types/verification';
import toast from 'react-hot-toast';

interface VerificationState {
  currentVerification: VerificationResponse | null;
  history: VerificationResponse[];
  stats: VerificationStats | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: VerificationState = {
  currentVerification: null,
  history: [],
  stats: null,
  isLoading: false,
  error: null,
};

export const verifyMessage = createAsyncThunk(
  'verification/verify',
  async (data: VerificationRequest, { rejectWithValue }) => {
    try {
      return await verificationService.verifyMessage(data);
    } catch (error: any) {
      const message = error.response?.data?.message || 'Verification failed';
      toast.error(message);
      return rejectWithValue(message);
    }
  }
);

export const fetchHistory = createAsyncThunk(
  'verification/fetchHistory',
  async ({ limit, offset }: { limit?: number; offset?: number }, { rejectWithValue }) => {
    try {
      return await verificationService.getHistory(limit, offset);
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch history');
    }
  }
);

export const fetchStats = createAsyncThunk(
  'verification/fetchStats',
  async (_, { rejectWithValue }) => {
    try {
      return await verificationService.getStats();
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch stats');
    }
  }
);

const verificationSlice = createSlice({
  name: 'verification',
  initialState,
  reducers: {
    clearCurrentVerification: (state) => {
      state.currentVerification = null;
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Verify Message
      .addCase(verifyMessage.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(verifyMessage.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentVerification = action.payload;
      })
      .addCase(verifyMessage.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch History
      .addCase(fetchHistory.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(fetchHistory.fulfilled, (state, action) => {
        state.isLoading = false;
        state.history = action.payload;
      })
      .addCase(fetchHistory.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch Stats
      .addCase(fetchStats.fulfilled, (state, action) => {
        state.stats = action.payload;
      });
  },
});

export const { clearCurrentVerification, clearError } = verificationSlice.actions;
export default verificationSlice.reducer;

