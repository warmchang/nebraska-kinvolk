import IconButton from '@mui/material/IconButton';
import Popover from '@mui/material/Popover';
import { styled } from '@mui/material/styles';
import React from 'react';
import { TwitterPicker } from 'react-color';

import ChannelAvatar from '../../Channels/ChannelAvatar';

const PREFIX = 'ColorPicker';

const classes = {
  iconButton: `${PREFIX}-iconButton`,
};

const Root = styled('div')({
  [`& .${classes.iconButton}`]: {
    padding: '0',
  },
});

export interface ColorPickerProps {
  /** Text color, applied to the children and used in the popover. */
  color: string;
  /** When a color is chosen. */
  onColorPicked: (color: { hex: string }) => void;
  /** Should the color picker be displayed initially? */
  initialOpen?: boolean;
  children: React.ReactElement<any>;
}

export default function ColorPicker(props: ColorPickerProps) {
  const [channelColor, setChannelColor] = React.useState(props.color);
  const [displayColorPicker, setDisplayColorPicker] = React.useState(props.initialOpen);
  const [anchorEl, setAnchorEl] = React.useState(null);
  const { onColorPicked } = props;

  function handleColorChange(color: { hex: string }) {
    setChannelColor(color.hex);
    onColorPicked(color);
  }

  function handleColorButtonClick(event: any) {
    setAnchorEl(event.currentTarget);
    setDisplayColorPicker(true);
  }

  function handleClose() {
    setAnchorEl(null);
    setDisplayColorPicker(false);
  }

  return (
    <Root>
      <IconButton
        className={classes.iconButton}
        onClick={handleColorButtonClick}
        data-testid="icon-button"
        size="large"
      >
        {props.children ? (
          <ChannelAvatar color={channelColor}>{props.children}</ChannelAvatar>
        ) : (
          <ChannelAvatar color={channelColor} />
        )}
      </IconButton>
      {displayColorPicker && anchorEl && (
        <Popover
          data-testid="popover"
          open={displayColorPicker}
          anchorEl={anchorEl}
          onClose={handleClose}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'center',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'center',
          }}
        >
          <TwitterPicker
            color={channelColor}
            onChangeComplete={handleColorChange}
            triangle="hide"
          />
        </Popover>
      )}
    </Root>
  );
}
