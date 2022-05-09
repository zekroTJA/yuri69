import styled from 'styled-components';
import { Disableable, Styled } from '../props';

type Props = Styled &
  Disableable & {
    min: number;
    max: number;
    value: number;
    onChange: (v: number) => void;
  };

const Label = styled.div`
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bolder;
  font-size: 1.2em;
  opacity: 0;

  transition: opacity 0.2s ease;
`;

const SliderContainer = styled.div<{ backgroundWidth: string }>`
  width: 100%;
  height: 100%;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: ${(p) => p.backgroundWidth};
    height: 100%;
    background-color: rgba(0 0 0 / 25%);
  }

  &:hover > ${Label} {
    opacity: 1;
  }
`;

const StyledInput = styled.input`
  width: 100%;

  -webkit-appearance: none;
  appearance: none;
  width: 100%;
  height: 4em;
  background-color: transparent;
  outline: none;
  opacity: 0.7;
  -webkit-transition: 0.2s;
  transition: opacity 0.2s;

  &::-webkit-slider-thumb {
    -webkit-appearance: none;
    background-color: transparent;
    opacity: 0;
  }

  &::-moz-range-thumb {
    background-color: transparent;
    opacity: 0;
  }
`;

export const Slider: React.FC<Props> = ({
  min,
  max,
  value = 0,
  onChange,
  disabled = false,
  ...props
}) => {
  const backgroundWidth = Math.floor(((value - min) * 100) / (max - min)).toString() + '%';
  return (
    <SliderContainer backgroundWidth={backgroundWidth} {...props}>
      <Label>
        <span>{value}</span>
      </Label>
      <StyledInput
        disabled={disabled}
        type="range"
        min={min}
        max={max}
        value={value}
        onInput={(e) => onChange(parseInt(e.currentTarget.value))}
      />
    </SliderContainer>
  );
};
