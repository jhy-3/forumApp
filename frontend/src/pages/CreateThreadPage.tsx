import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Typography,
  Card,
  CardContent,
  TextField,
  Button,
  Box,
  Alert,
} from '@mui/material';
import { threadAPI, forumAPI } from '../api/api';
import { Forum } from '../types';

const CreateThreadPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const [forum, setForum] = useState<Forum | null>(null);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    loadForum();
  }, [slug]);

  const loadForum = async () => {
    try {
      const response = await forumAPI.getAll();
      const foundForum = response.data.find((f) => f.slug === slug);
      if (foundForum) {
        setForum(foundForum);
      }
    } catch (error) {
      console.error('Failed to load forum:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!forum) return;

    setError('');
    setSubmitting(true);

    try {
      const response = await threadAPI.create({
        forumId: forum.id,
        title,
        content,
      });
      navigate(`/threads/${response.data.id}`);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to create thread');
      setSubmitting(false);
    }
  };

  if (!forum) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
        <Typography>Loading...</Typography>
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Create New Thread in {forum.name}
      </Typography>
      <Card>
        <CardContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          <form onSubmit={handleSubmit}>
            <TextField
              label="Thread Title"
              fullWidth
              margin="normal"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
            <TextField
              label="Content"
              multiline
              rows={10}
              fullWidth
              margin="normal"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              required
            />
            <Box display="flex" gap={2} mt={2}>
              <Button
                type="submit"
                variant="contained"
                disabled={submitting}
              >
                {submitting ? 'Creating...' : 'Create Thread'}
              </Button>
              <Button
                variant="outlined"
                onClick={() => navigate(`/forums/${slug}`)}
              >
                Cancel
              </Button>
            </Box>
          </form>
        </CardContent>
      </Card>
    </Box>
  );
};

export default CreateThreadPage;

