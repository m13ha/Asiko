import API_CONFIG from './ApiConfig';

/**
 * A simple CORS proxy service for development
 */
class CorsProxy {
  private baseUrl = API_CONFIG.API_URL;
  
  /**
   * Make a fetch request through a CORS proxy
   */
  async fetch(url: string, options: RequestInit = {}): Promise<Response> {
    // Add CORS headers to all requests
    const headers = {
      ...options.headers,
      'Origin': window.location.origin,
    };
    
    // Set mode to 'cors' to enable CORS
    const corsOptions: RequestInit = {
      ...options,
      headers,
      mode: 'cors',
    };
    
    try {
      return await fetch(url, corsOptions);
    } catch (error) {
      console.error('CORS proxy fetch error:', error);
      throw error;
    }
  }
}

const corsProxy = new CorsProxy();
export default corsProxy;