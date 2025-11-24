import React, { useState, useEffect } from 'react';
import {
  Typography,
  Card,
  CardContent,
  TextField,
  Button,
  Box,
  Alert,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import { userAPI } from '../api/api';

const ProfilePage: React.FC = () => {
  const { user } = useAuth();
  const [avatarUrl, setAvatarUrl] = useState('');
  const [aboutText, setAboutText] = useState('');
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (user) {
      setAvatarUrl(user.avatarUrl || '');
      setAboutText(user.aboutText || '');
    }
  }, [user]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setSubmitting(true);

    try {
      await userAPI.updateProfile({
        avatarUrl: avatarUrl || undefined,
        aboutText: aboutText || undefined,
      });
      setSuccess('Profile updated successfully!');
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to update profile');
    } finally {
      setSubmitting(false);
    }
  };

  if (!user) {
    return null;
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Profile Settings
      </Typography>
      <Card>
        <CardContent>
          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              {success}
            </Alert>
          )}
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          <Box mb={3}>
            <Typography variant="subtitle1" gutterBottom>
              Username: <strong>{user.username}</strong>
            </Typography>
            <Typography variant="subtitle1" gutterBottom>
              Email: <strong>{user.email}</strong>
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Member since: {new Date(user.createdAt).toLocaleDateString()}
            </Typography>
          </Box>
          <form onSubmit={handleSubmit}>
            <TextField
              label="Avatar URL"
              fullWidth
              margin="normal"
              value={avatarUrl}
              onChange={(e) => setAvatarUrl(e.target.value)}
              helperText="URL to your profile picture"
            />
            <TextField
              label="About Me"
              multiline
              rows={4}
              fullWidth
              margin="normal"
              value={aboutText}
              onChange={(e) => setAboutText(e.target.value)}
              helperText="Tell us about yourself"
            />
            <Button
              type="submit"
              variant="contained"
              sx={{ mt: 2 }}
              disabled={submitting}
            >
              {submitting ? 'Updating...' : 'Update Profile'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </Box>
  );
};

export default ProfilePage;

