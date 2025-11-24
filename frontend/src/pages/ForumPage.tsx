import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  Typography,
  Card,
  CardContent,
  Box,
  Button,
  CircularProgress,
  Chip,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { Thread } from '../types';
import { forumAPI } from '../api/api';
import { useAuth } from '../contexts/AuthContext';

const ForumPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const [threads, setThreads] = useState<Thread[]>([]);
  const [loading, setLoading] = useState(true);
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    if (slug) {
      loadThreads();
    }
  }, [slug]);

  const loadThreads = async () => {
    try {
      if (slug) {
        const response = await forumAPI.getThreads(slug);
        setThreads(response.data.threads || []);
      }
    } catch (error) {
      console.error('Failed to load threads:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Threads</Typography>
        {isAuthenticated && (
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            component={Link}
            to={`/forums/${slug}/create-thread`}
          >
            New Thread
          </Button>
        )}
      </Box>

      {threads.length === 0 ? (
        <Typography variant="body1" color="text.secondary">
          No threads yet. Be the first to create one!
        </Typography>
      ) : (
        threads.map((thread) => (
          <Card key={thread.id} sx={{ mb: 2 }}>
            <CardContent>
              <Box display="flex" alignItems="center" gap={1} mb={1}>
                {thread.isPinned && <Chip label="Pinned" size="small" color="primary" />}
                {thread.isLocked && <Chip label="Locked" size="small" />}
              </Box>
              <Typography
                variant="h6"
                component={Link}
                to={`/threads/${thread.id}`}
                sx={{ textDecoration: 'none', color: 'primary.main' }}
              >
                {thread.title}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                By {thread.author?.username} • {thread.replyCount} replies • {thread.viewCount}{' '}
                views
              </Typography>
            </CardContent>
          </Card>
        ))
      )}
    </Box>
  );
};

export default ForumPage;

