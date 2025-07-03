'use client';

import { useState } from 'react';
import { toast } from 'react-hot-toast';

interface BuildServiceReturn {
  building: boolean;
  deploying: boolean;
  buildProject: (projectId: string) => Promise<void>;
  deployProject: (projectId: string) => Promise<void>;
}

export function useBuildService(): BuildServiceReturn {
  const [building, setBuilding] = useState(false);
  const [deploying, setDeploying] = useState(false);

  const buildProject = async (projectId: string): Promise<void> => {
    if (building) return;

    setBuilding(true);

    try {
      const response = await fetch(`/api/projects/${projectId}/build`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error('Failed to start build');
      }

      const result = await response.json();

      if (result.success) {
        toast.success('Build iniciado com sucesso!');
      } else {
        throw new Error(result.error?.message || 'Build failed');
      }
    } catch (error) {
      console.error('Build error:', error);
      toast.error('Erro ao iniciar build');
      throw error;
    } finally {
      setBuilding(false);
    }
  };

  const deployProject = async (projectId: string): Promise<void> => {
    if (deploying) return;

    setDeploying(true);

    try {
      const response = await fetch(`/api/projects/${projectId}/deploy`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error('Failed to start deployment');
      }

      const result = await response.json();

      if (result.success) {
        toast.success('Deploy iniciado com sucesso!');
      } else {
        throw new Error(result.error?.message || 'Deploy failed');
      }
    } catch (error) {
      console.error('Deploy error:', error);
      toast.error('Erro ao iniciar deploy');
      throw error;
    } finally {
      setDeploying(false);
    }
  };

  return {
    building,
    deploying,
    buildProject,
    deployProject,
  };
}
