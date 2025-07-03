'use client';

import { cn } from '@/lib/utils';
import React from 'react';

interface SidebarProps {
  children: React.ReactNode;
  className?: string;
}

export function Sidebar({ children, className }: SidebarProps) {
  return (
    <aside className={cn('bg-white border-r border-gray-200 overflow-y-auto', className)}>
      {children}
    </aside>
  );
}

interface ComponentLibraryProps {
  className?: string;
}

export function ComponentLibrary({ className }: ComponentLibraryProps) {
  const components = [
    { id: 'text', name: 'Texto', icon: '📝' },
    { id: 'heading', name: 'Título', icon: '📰' },
    { id: 'image', name: 'Imagem', icon: '🖼️' },
    { id: 'button', name: 'Botão', icon: '🔘' },
    { id: 'container', name: 'Container', icon: '📦' },
    { id: 'grid', name: 'Grid', icon: '⚏' },
    { id: 'form', name: 'Formulário', icon: '📋' },
    { id: 'video', name: 'Vídeo', icon: '🎥' },
  ];

  return (
    <div className={cn('p-4', className)}>
      <h3 className="text-sm font-semibold text-gray-900 mb-4">Componentes</h3>

      <div className="space-y-2">
        {components.map((component) => (
          <div
            key={component.id}
            draggable
            className="flex items-center p-3 bg-gray-50 rounded-lg cursor-move hover:bg-gray-100 transition-colors"
          >
            <span className="text-lg mr-3">{component.icon}</span>
            <span className="text-sm font-medium text-gray-700">{component.name}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

interface CanvasProps {
  project: any;
  selectedElement: string | null;
  onSelectElement: (elementId: string | null) => void;
  className?: string;
}

export const Canvas = React.forwardRef<HTMLDivElement, CanvasProps>(
  ({ project, selectedElement, onSelectElement, className }, ref) => {
    return (
      <div
        ref={ref}
        className={cn('flex-1 bg-gray-100 overflow-auto p-4', className)}
      >
        <div className="max-w-6xl mx-auto bg-white min-h-screen shadow-lg">
          {project?.sections?.map((section: any, index: number) => (
            <div
              key={section.id}
              className={cn(
                'border-2 border-transparent hover:border-blue-300 transition-colors',
                selectedElement === section.id && 'border-blue-500'
              )}
              onClick={() => onSelectElement(section.id)}
            >
              <div className="p-4">
                <h3 className="text-lg font-semibold mb-2">{section.name}</h3>
                <div className="min-h-[100px] bg-gray-50 rounded border-2 border-dashed border-gray-300 flex items-center justify-center">
                  <span className="text-gray-500">Seção {index + 1}</span>
                </div>
              </div>
            </div>
          )) || (
            <div className="flex items-center justify-center h-64">
              <div className="text-center">
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  Comece a criar seu site
                </h3>
                <p className="text-gray-600">
                  Arraste componentes da biblioteca para começar
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    );
  }
);

Canvas.displayName = 'Canvas';

interface PropertyPanelProps {
  elementId: string;
  project: any;
  onUpdate: (updates: any) => void;
  className?: string;
}

export function PropertyPanel({ elementId, project, onUpdate, className }: PropertyPanelProps) {
  return (
    <div className={cn('p-4 border-b border-gray-200', className)}>
      <h3 className="text-sm font-semibold text-gray-900 mb-4">Propriedades</h3>

      <div className="space-y-4">
        <div>
          <label className="block text-xs font-medium text-gray-700 mb-1">
            Nome
          </label>
          <input
            type="text"
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
            placeholder="Nome do elemento"
          />
        </div>

        <div>
          <label className="block text-xs font-medium text-gray-700 mb-1">
            Texto
          </label>
          <textarea
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm resize-none"
            rows={3}
            placeholder="Conteúdo do elemento"
          />
        </div>
      </div>
    </div>
  );
}

interface StylesPanelProps {
  elementId: string;
  project: any;
  onUpdate: (updates: any) => void;
  className?: string;
}

export function StylesPanel({ elementId, project, onUpdate, className }: StylesPanelProps) {
  return (
    <div className={cn('p-4', className)}>
      <h3 className="text-sm font-semibold text-gray-900 mb-4">Estilos</h3>

      <div className="space-y-4">
        <div>
          <label className="block text-xs font-medium text-gray-700 mb-1">
            Cor de Fundo
          </label>
          <input
            type="color"
            className="w-full h-8 border border-gray-300 rounded-md"
          />
        </div>

        <div>
          <label className="block text-xs font-medium text-gray-700 mb-1">
            Cor do Texto
          </label>
          <input
            type="color"
            className="w-full h-8 border border-gray-300 rounded-md"
          />
        </div>

        <div>
          <label className="block text-xs font-medium text-gray-700 mb-1">
            Tamanho da Fonte
          </label>
          <select className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm">
            <option value="12">12px</option>
            <option value="14">14px</option>
            <option value="16">16px</option>
            <option value="18">18px</option>
            <option value="20">20px</option>
            <option value="24">24px</option>
          </select>
        </div>
      </div>
    </div>
  );
}

interface PreviewPanelProps {
  project: any;
  className?: string;
}

export function PreviewPanel({ project, className }: PreviewPanelProps) {
  return (
    <div className={cn('flex-1 bg-white', className)}>
      <iframe
        src={`/preview/${project?.id}`}
        className="w-full h-full border-0"
        title="Preview"
      />
    </div>
  );
}

interface CodePanelProps {
  project: any;
  className?: string;
}

export function CodePanel({ project, className }: CodePanelProps) {
  return (
    <div className={cn('flex-1 bg-gray-900 text-green-400 p-4 font-mono text-sm', className)}>
      <pre className="whitespace-pre-wrap">
        {JSON.stringify(project, null, 2)}
      </pre>
    </div>
  );
}

interface BuildProgressModalProps {
  progress: {
    status: string;
    progress: number;
    message: string;
    logs: string[];
  };
  onClose: () => void;
}

export function BuildProgressModal({ progress, onClose }: BuildProgressModalProps) {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-96 max-w-full">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">Build em Progresso</h3>
          {progress.status === 'success' || progress.status === 'error' ? (
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700"
            >
              ✕
            </button>
          ) : null}
        </div>

        <div className="space-y-4">
          <div>
            <div className="flex justify-between text-sm text-gray-600 mb-1">
              <span>{progress.message}</span>
              <span>{progress.progress}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${progress.progress}%` }}
              />
            </div>
          </div>

          {progress.logs.length > 0 && (
            <div className="bg-gray-100 rounded p-3 max-h-32 overflow-y-auto">
              <div className="text-xs text-gray-600 space-y-1">
                {progress.logs.map((log, index) => (
                  <div key={index}>{log}</div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
