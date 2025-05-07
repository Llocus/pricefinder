import { useState, useEffect, useCallback } from 'react';

export default function useWindowDimensions() {

  const hasWindow = typeof window !== 'undefined';

  const getWindowDimensions = useCallback(() => {
    const width = hasWindow ? window.innerWidth : 0;
    const height = hasWindow ? window.innerHeight : 0;
    return {
      width,
      height,
    };
 }, [hasWindow]);

  const [windowDimensions, setWindowDimensions] = useState(getWindowDimensions());

  useEffect(() => {
    if (hasWindow) {
      const handleResize = () => {
        setWindowDimensions(getWindowDimensions());
      }

      window.addEventListener('resize', handleResize);
      return () => window.removeEventListener('resize', handleResize);
    }
  }, [getWindowDimensions, hasWindow]);

  return windowDimensions;
}