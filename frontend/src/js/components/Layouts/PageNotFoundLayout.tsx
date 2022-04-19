import { Grid, Typography } from '@material-ui/core';
import pageNotFoundGraphic from '../../../img/404.svg';

export default function PageNotFoundLayout() {
  return (
    <Grid
      container
      spacing={0}
      direction="column"
      alignItems="center"
      justify="center"
      style={{ minHeight: '100vh', textAlign: 'center' }}
    >
      <Grid item xs={6}>
        <img src={pageNotFoundGraphic} alt="page not found 404" style={{ maxWidth: '100%' }} />
        <Typography variant="h5" paragraph>
          Whoops! This page doesn't exist
        </Typography>
        <Typography variant="h6" paragraph>
          Not to worry, head back to <a href="/">home</a>!
        </Typography>{' '}
      </Grid>
    </Grid>
  );
}
