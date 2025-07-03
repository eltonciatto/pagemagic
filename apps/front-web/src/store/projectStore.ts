import { Project } from '@/types';
import { create } from 'zustand';

interface ProjectState {
  project: Project | null;
  projects: Project[];
  loading: boolean;
  saving: boolean;

  // Actions
  setProject: (project: Project) => void;
  setProjects: (projects: Project[]) => void;
  updateProject: (updates: Partial<Project>) => void;
  addProject: (project: Project) => void;
  removeProject: (projectId: string) => void;
  setLoading: (loading: boolean) => void;
  setSaving: (saving: boolean) => void;
  saveProject: () => Promise<void>;
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  project: null,
  projects: [],
  loading: false,
  saving: false,

  setProject: (project) => {
    set({ project });
  },

  setProjects: (projects) => {
    set({ projects });
  },

  updateProject: (updates) => {
    const { project } = get();
    if (project) {
      const updatedProject = { ...project, ...updates };
      set({ project: updatedProject });

      // Also update in projects list if present
      const { projects } = get();
      const projectIndex = projects.findIndex(p => p.id === project.id);
      if (projectIndex !== -1) {
        const updatedProjects = [...projects];
        updatedProjects[projectIndex] = updatedProject;
        set({ projects: updatedProjects });
      }
    }
  },

  addProject: (project) => {
    const { projects } = get();
    set({ projects: [project, ...projects] });
  },

  removeProject: (projectId) => {
    const { projects } = get();
    set({ projects: projects.filter(p => p.id !== projectId) });

    const { project } = get();
    if (project && project.id === projectId) {
      set({ project: null });
    }
  },

  setLoading: (loading) => {
    set({ loading });
  },

  setSaving: (saving) => {
    set({ saving });
  },

  saveProject: async () => {
    const { project, setSaving } = get();
    if (!project) return;

    setSaving(true);

    try {
      const response = await fetch(`/api/projects/${project.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(project),
      });

      if (!response.ok) {
        throw new Error('Failed to save project');
      }

      const updatedProject = await response.json();
      set({ project: updatedProject });

      console.log('Project saved successfully');
    } catch (error) {
      console.error('Error saving project:', error);
      throw error;
    } finally {
      setSaving(false);
    }
  },
}));
