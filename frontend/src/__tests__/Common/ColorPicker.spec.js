import { fireEvent, render } from '@testing-library/react';
import React from 'react';
import ColorPicker from '../../components/Common/ColorPicker';

describe('Color Picker', () => {
  it('should render popover on Icon Click', () => {
    const { getByTestId } = render(<ColorPicker color={'#fff'} />);
    fireEvent.click(getByTestId('icon-button'));
    expect(getByTestId('popover')).toBeTruthy();
  });
});
