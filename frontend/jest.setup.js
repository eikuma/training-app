// jest.setup.js
import '@testing-library/jest-dom'; // Changed from /extend-expect

// Mock next/router for components using next/link or useRouter
jest.mock('next/router', () => ({
  useRouter() {
    return {
      route: '/',
      pathname: '',
      query: '',
      asPath: '',
      push: jest.fn(),
      events: {
        on: jest.fn(),
        off: jest.fn()
      },
      beforePopState: jest.fn(() => null),
      prefetch: jest.fn(() => null),
      replace: jest.fn(),
    };
  },
}));

// Mock next/link
jest.mock('next/link', () => {
  const React = require('react');
  // eslint-disable-next-line react/display-name
  return ({ children, href, passHref, ...rest }) => { // Explicitly destructure passHref
    // If children is a simple string (like the "Fitness Tracker" title link)
    if (typeof children === 'string') {
      // Do not spread passHref to the DOM element
      return <a href={href} {...rest}>{children}</a>;
    }
    // If children is a single React element (like a Button)
    // Ensure that we only attempt to clone if children is a valid element.
    if (React.isValidElement(children)) {
       // Do not spread passHref to the DOM element if it's a DOM element,
       // but if it's a custom component that expects it, this might need adjustment.
       // For MUI Button, it doesn't need passHref.
       return React.cloneElement(children, { href, ...rest });
    }
    // Fallback for other cases (e.g., multiple children, though Link usually takes one)
    return <a href={href} {...rest}>{children}</a>;
  };
});
