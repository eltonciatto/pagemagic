'use client';

import { DragDropContext } from '@hello-pangea/dnd';
import {
  ArrowUturnLeftIcon,
  ArrowUturnRightIcon,
  CodeBracketIcon,
  DocumentArrowDownIcon,
  EyeIcon,
  PaintBrushIcon,
  PauseIcon,
  PlayIcon
} from '@heroicons/react/24/outline';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useRef, useState } from 'react';
import { toast } from 'react-hot-toast';

import {
  BuildProgressModal,
  Canvas,
  CodePanel,
  ComponentLibrary,
  PreviewPanel,
  PropertyPanel,
  Sidebar,
  StylesPanel
} from '@/components/editor';
import { Button } from '@/components/ui/Button';
import { useBuildService } from '@/hooks/useBuildService';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useProjectStore } from '@/store/projectStore';
import { Component, Project } from '@/types';

type EditorMode = 'design' | 'preview' | 'code';

interface BuildProgress {
  status: 'idle' | 'building' | 'deploying' | 'success' | 'error';
  progress: number;
  message: string;
  logs: string[];
}

export default function EditorPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const { project, loading, updateProject, saveProject } = useProjectStore();
  const { buildProject, deployProject } = useBuildService();
  const { socket, connected } = useWebSocket();

  const [mode, setMode] = useState<EditorMode>('design');
  const [selectedElement, setSelectedElement] = useState<string | null>(null);
  const [showComponentLibrary, setShowComponentLibrary] = useState(true);
  const [showPropertyPanel, setShowPropertyPanel] = useState(true);
  const [buildProgress, setBuildProgress] = useState<BuildProgress>({
    status: 'idle',
    progress: 0,
    message: '',
    logs: []
  });
  const [undoStack, setUndoStack] = useState<Project[]>([]);
  const [redoStack, setRedoStack] = useState<Project[]>([]);

  const canvasRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!project && projectId) {
      // Load project from API
      loadProject(projectId);
    }
  }, [projectId, project]);

  useEffect(() => {
    if (socket && connected) {
      socket.on('buildProgress', handleBuildProgress);
      socket.on('buildComplete', handleBuildComplete);
      socket.on('buildError', handleBuildError);

      return () => {
        socket.off('buildProgress');
        socket.off('buildComplete');
        socket.off('buildError');
      };
    }
  }, [socket, connected]);

  const loadProject = async (id: string) => {
    try {
      const response = await fetch(`/api/projects/${id}`);
      if (!response.ok) throw new Error('Failed to load project');

      const projectData = await response.json();
      updateProject(projectData);
    } catch (error) {
      console.error('Error loading project:', error);
      toast.error('Falha ao carregar projeto');
      router.push('/dashboard');
    }
  };

  const handleDragEnd = (result: any) => {
    if (!result.destination || !project) return;

    const { source, destination, draggableId } = result;

    // Create undo point
    addToUndoStack();

    if (source.droppableId === 'component-library') {
      // Adding new component from library
      const newComponent: Component = {
        id: `component_${Date.now()}`,
        type: draggableId,
        props: {},
        styles: {},
        children: []
      };

      const updatedProject = addComponentToSection(project, destination.droppableId, newComponent, destination.index);
      updateProject(updatedProject);
    } else {
      // Moving existing component
      const updatedProject = moveComponent(project, source, destination);
      updateProject(updatedProject);
    }

    // Auto-save after 2 seconds
    setTimeout(() => {
      saveProject();
    }, 2000);
  };

  const addComponentToSection = (project: Project, sectionId: string, component: Component, index: number): Project => {
    return {
      ...project,
      sections: project.sections.map(section => {
        if (section.id === sectionId) {
          const newComponents = [...section.components];
          newComponents.splice(index, 0, component);
          return { ...section, components: newComponents };
        }
        return section;
      })
    };
  };

  const moveComponent = (project: Project, source: any, destination: any): Project => {
    // Implementation for moving components between sections
    // This would handle the complex logic of moving components
    return project;
  };

  const addToUndoStack = () => {
    if (!project) return;

    setUndoStack(prev => [...prev.slice(-19), project]);
    setRedoStack([]);
  };

  const undo = () => {
    if (undoStack.length === 0 || !project) return;

    const previousState = undoStack[undoStack.length - 1];
    setRedoStack(prev => [...prev, project]);
    setUndoStack(prev => prev.slice(0, -1));
    updateProject(previousState);
  };

  const redo = () => {
    if (redoStack.length === 0) return;

    const nextState = redoStack[redoStack.length - 1];
    setUndoStack(prev => [...prev, project!]);
    setRedoStack(prev => prev.slice(0, -1));
    updateProject(nextState);
  };

  const handleBuild = async () => {
    if (!project) return;

    setBuildProgress({
      status: 'building',
      progress: 0,
      message: 'Iniciando build...',
      logs: []
    });

    try {
      await buildProject(project.id);
    } catch (error) {
      console.error('Build error:', error);
      setBuildProgress(prev => ({
        ...prev,
        status: 'error',
        message: 'Erro no build'
      }));
      toast.error('Erro no build do projeto');
    }
  };

  const handleDeploy = async () => {
    if (!project) return;

    setBuildProgress({
      status: 'deploying',
      progress: 50,
      message: 'Fazendo deploy...',
      logs: []
    });

    try {
      await deployProject(project.id);
    } catch (error) {
      console.error('Deploy error:', error);
      setBuildProgress(prev => ({
        ...prev,
        status: 'error',
        message: 'Erro no deploy'
      }));
      toast.error('Erro no deploy do projeto');
    }
  };

  const handleBuildProgress = (data: any) => {
    setBuildProgress(prev => ({
      ...prev,
      progress: data.progress,
      message: data.message,
      logs: [...prev.logs, data.log]
    }));
  };

  const handleBuildComplete = (data: any) => {
    setBuildProgress(prev => ({
      ...prev,
      status: 'success',
      progress: 100,
      message: 'Build concluído com sucesso!'
    }));
    toast.success('Build concluído!');
  };

  const handleBuildError = (data: any) => {
    setBuildProgress(prev => ({
      ...prev,
      status: 'error',
      message: data.message || 'Erro no build'
    }));
    toast.error('Erro no build');
  };

  if (loading || !project) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-2 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <h1 className="text-lg font-semibold text-gray-900">{project.name}</h1>
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={undo}
              disabled={undoStack.length === 0}
            >
              <ArrowUturnLeftIcon className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={redo}
              disabled={redoStack.length === 0}
            >
              <ArrowUturnRightIcon className="h-4 w-4" />
            </Button>
          </div>
        </div>

        <div className="flex items-center space-x-2">
          {/* Mode Switcher */}
          <div className="flex bg-gray-100 rounded-lg p-1">
            <button
              onClick={() => setMode('design')}
              className={`px-3 py-1 rounded-md text-sm font-medium ${
                mode === 'design'
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              <PaintBrushIcon className="h-4 w-4 inline mr-1" />
              Design
            </button>
            <button
              onClick={() => setMode('preview')}
              className={`px-3 py-1 rounded-md text-sm font-medium ${
                mode === 'preview'
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              <EyeIcon className="h-4 w-4 inline mr-1" />
              Preview
            </button>
            <button
              onClick={() => setMode('code')}
              className={`px-3 py-1 rounded-md text-sm font-medium ${
                mode === 'code'
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              <CodeBracketIcon className="h-4 w-4 inline mr-1" />
              Code
            </button>
          </div>

          {/* Actions */}
          <Button onClick={handleBuild} disabled={buildProgress.status === 'building'}>
            {buildProgress.status === 'building' ? (
              <PauseIcon className="h-4 w-4 mr-2" />
            ) : (
              <PlayIcon className="h-4 w-4 mr-2" />
            )}
            Build
          </Button>

          <Button onClick={handleDeploy} disabled={buildProgress.status !== 'idle'}>
            <DocumentArrowDownIcon className="h-4 w-4 mr-2" />
            Deploy
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        <DragDropContext onDragEnd={handleDragEnd}>
          {/* Left Sidebar */}
          {showComponentLibrary && (
            <Sidebar className="w-64 border-r border-gray-200">
              <ComponentLibrary />
            </Sidebar>
          )}

          {/* Main Canvas Area */}
          <div className="flex-1 flex flex-col">
            {mode === 'design' && (
              <Canvas
                ref={canvasRef}
                project={project}
                selectedElement={selectedElement}
                onSelectElement={setSelectedElement}
              />
            )}

            {mode === 'preview' && (
              <PreviewPanel project={project} />
            )}

            {mode === 'code' && (
              <CodePanel project={project} />
            )}
          </div>

          {/* Right Panel */}
          {showPropertyPanel && mode === 'design' && (
            <div className="w-80 border-l border-gray-200 bg-white">
              {selectedElement ? (
                <>
                  <PropertyPanel
                    elementId={selectedElement}
                    project={project}
                    onUpdate={updateProject}
                  />
                  <StylesPanel
                    elementId={selectedElement}
                    project={project}
                    onUpdate={updateProject}
                  />
                </>
              ) : (
                <div className="p-4 text-center text-gray-500">
                  Selecione um elemento para editar suas propriedades
                </div>
              )}
            </div>
          )}
        </DragDropContext>
      </div>

      {/* Build Progress Modal */}
      {buildProgress.status !== 'idle' && (
        <BuildProgressModal
          progress={buildProgress}
          onClose={() => setBuildProgress({ status: 'idle', progress: 0, message: '', logs: [] })}
        />
      )}
    </div>
  );
}
