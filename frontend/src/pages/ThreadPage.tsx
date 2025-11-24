import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Typography,
  Card,
  CardContent,
  Box,
  TextField,
  Button,
  CircularProgress,
  Divider,
  Alert,
} from '@mui/material';
import { Thread, Post } from '../types';
import { threadAPI, postAPI } from '../api/api';
import { useAuth } from '../contexts/AuthContext';

const ThreadPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [thread, setThread] = useState<Thread | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [replyContent, setReplyContent] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    if (id) {
      loadThread();
      loadPosts();
    }
  }, [id]);

  const loadThread = async () => {
    try {
      if (id) {
        const response = await threadAPI.getById(parseInt(id));
        setThread(response.data);
      }
    } catch (error) {
      console.error('Failed to load thread:', error);
    }
  };

  const loadPosts = async () => {
    try {
      if (id) {
        const response = await threadAPI.getPosts(parseInt(id));
        setPosts(response.data.posts || []);
      }
    } catch (error) {
      console.error('Failed to load posts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!replyContent.trim() || !id) return;

    setError('');
    setSubmitting(true);

    try {
      await postAPI.create(parseInt(id), { content: replyContent });
      setReplyContent('');
      loadPosts();
      loadThread(); // Reload to update reply count
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to post reply');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
        <CircularProgress />
      </Box>
    );
  }

  if (!thread) {
    return <Typography>Thread not found</Typography>;
  }

  return (
    <Box>
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h4" gutterBottom>
            {thread.title}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            By {thread.author?.username} • {new Date(thread.createdAt).toLocaleString()}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {thread.viewCount} views • {thread.replyCount} replies
          </Typography>
        </CardContent>
      </Card>

      <Typography variant="h6" gutterBottom>
        Replies
      </Typography>

      {posts.map((post, index) => (
        <Card key={post.id} sx={{ mb: 2 }}>
          <CardContent>
            <Box display="flex" justifyContent="space-between" alignItems="start" mb={1}>
              <Typography variant="subtitle2" color="primary">
                {post.author?.username}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                #{post.postNumber} • {new Date(post.createdAt).toLocaleString()}
              </Typography>
            </Box>
            <Divider sx={{ my: 1 }} />
            <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
              {post.content}
            </Typography>
          </CardContent>
        </Card>
      ))}

      {isAuthenticated && !thread.isLocked && (
        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Post Reply
            </Typography>
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}
            <form onSubmit={handleReply}>
              <TextField
                label="Your reply"
                multiline
                rows={4}
                fullWidth
                value={replyContent}
                onChange={(e) => setReplyContent(e.target.value)}
                required
              />
              <Button
                type="submit"
                variant="contained"
                sx={{ mt: 2 }}
                disabled={submitting}
              >
                {submitting ? 'Posting...' : 'Post Reply'}
              </Button>
            </form>
          </CardContent>
        </Card>
      )}

      {thread.isLocked && (
        <Alert severity="info" sx={{ mt: 3 }}>
          This thread is locked. No new replies can be posted.
        </Alert>
      )}

      {!isAuthenticated && (
        <Alert severity="info" sx={{ mt: 3 }}>
          Please log in to post a reply.
        </Alert>
      )}
    </Box>
  );
};

export default ThreadPage;

