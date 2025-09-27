/** @jsxRuntime classic */
/** @jsx jsx */
import { css, jsx } from "@emotion/react";

export const IconContainer = ({ color, children }) => (
  <span
    css={css`
      svg {
        color: ${color};
        text-align: center;
        vertical-align: -0.125em;
      }
    `}
  >
    {children}
  </span>
);
