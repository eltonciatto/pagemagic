import { useState, useEffect, createContext, useContext } from 'react';

export interface User {
  id: string;
  email: string;
  name?: string;
  avatar?: string;
  role: string;
  subscription?: {
    plan: string;
    status: string;
    current_period_end: string;
  };
  created_at: string;
  updated_at: string;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  register: (email: string, password: string, name: string) => Promise<void>;
  updateProfile: (data: Partial<User>) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Verificar se há token armazenado
    const token = localStorage.getItem('authToken');
    if (token) {
      // Validar token e carregar dados do usuário
      loadUserProfile();
    } else {
      setLoading(false);
    }
  }, []);

  const loadUserProfile = async () => {
    try {
      // Simular carregamento do perfil
      // Em implementação real, faria chamada para a API
      const mockUser: User = {
        id: 'user-1',
        email: 'user@example.com',
        name: 'Usuário Teste',
        role: 'user',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      };
      setUser(mockUser);
    } catch (error) {
      console.error('Erro ao carregar perfil:', error);
      localStorage.removeItem('authToken');
      localStorage.removeItem('userId');
    } finally {
      setLoading(false);
    }
  };

  const login = async (email: string, password: string) => {
    try {
      // Simular login
      // Em implementação real, faria chamada para authApi.login
      const mockToken = 'mock-jwt-token';
      const mockUser: User = {
        id: 'user-1',
        email,
        name: 'Usuário Teste',
        role: 'user',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      };

      localStorage.setItem('authToken', mockToken);
      localStorage.setItem('userId', mockUser.id);
      setUser(mockUser);
    } catch (error) {
      throw new Error('Credenciais inválidas');
    }
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('userId');
    setUser(null);
  };

  const register = async (email: string, password: string, name: string) => {
    try {
      // Simular registro
      // Em implementação real, faria chamada para authApi.register
      const mockToken = 'mock-jwt-token';
      const mockUser: User = {
        id: 'user-1',
        email,
        name,
        role: 'user',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      };

      localStorage.setItem('authToken', mockToken);
      localStorage.setItem('userId', mockUser.id);
      setUser(mockUser);
    } catch (error) {
      throw new Error('Erro ao criar conta');
    }
  };

  const updateProfile = async (data: Partial<User>) => {
    try {
      // Simular atualização de perfil
      if (user) {
        const updatedUser = { ...user, ...data, updated_at: new Date().toISOString() };
        setUser(updatedUser);
      }
    } catch (error) {
      throw new Error('Erro ao atualizar perfil');
    }
  };

  const value = {
    user,
    loading,
    login,
    logout,
    register,
    updateProfile
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth deve ser usado dentro de um AuthProvider');
  }
  return context;
}
