console.log('Desenvolvimento do Page Magic - Prompt Service');

export const logger = {
  info: (message: string) => console.log(`[INFO] ${new Date().toISOString()} - ${message}`),
  error: (message: string) => console.error(`[ERROR] ${new Date().toISOString()} - ${message}`),
  warn: (message: string) => console.warn(`[WARN] ${new Date().toISOString()} - ${message}`),
  debug: (message: string) => console.log(`[DEBUG] ${new Date().toISOString()} - ${message}`),
};
