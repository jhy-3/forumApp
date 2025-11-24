import React, { useEffect, useState } from 'react';
import {
  Typography,
  Card,
  CardContent,
  Grid,
  Box,
  CircularProgress,
} from '@mui/material';
import { Link } from 'react-router-dom';
import { Forum } from '../types';
import { forumAPI } from '../api/api';

const HomePage: React.FC = () => {
  const [forums, setForums] = useState<Forum[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadForums();
  }, []);

  const loadForums = async () => {
    try {
      const response = await forumAPI.getAll();
      setForums(response.data);
    } catch (error) {
      console.error('Failed to load forums:', error);
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
      <Typography variant="h4" gutterBottom>
        Forums
      </Typography>
      <Grid container spacing={2}>
        {forums.map((forum) => (
          <Grid item xs={12} key={forum.id}>
            <Card component={Link} to={`/forums/${forum.slug}`} sx={{ textDecoration: 'none' }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  {forum.name}
                </Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  {forum.description}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {forum.threadCount} threads • {forum.postCount} posts
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default HomePage;

