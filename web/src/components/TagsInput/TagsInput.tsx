import React, { useEffect, useState } from 'react';
import { Input } from '../Input';
import { Disableable } from '../props';

type Props = React.HTMLAttributes<HTMLInputElement> &
  Disableable & {
    tags?: string[];
    onTagsChange?: (v: string[]) => void;
  };

export const TagsInput: React.FC<Props> = ({ tags = [], onTagsChange = () => {}, ...props }) => {
  const [tagsValue, setTagsValue] = useState('');

  const _valueToTags = (v: string) =>
    v
      .split(',')
      .map((t) => t.trim())
      .filter((t) => !!t);

  const _tagsToValue = (t: string[]) => t.join(', ');

  useEffect(() => {
    setTagsValue(_tagsToValue(tags ?? []));
  }, [tags]);

  return (
    <Input
      value={tagsValue}
      onInput={(e) => setTagsValue(e.currentTarget.value)}
      onBlur={(e) => onTagsChange(_valueToTags(e.currentTarget.value))}
      {...props}></Input>
  );
};
