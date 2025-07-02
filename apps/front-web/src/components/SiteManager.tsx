'use client';

import React, { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { apiClient } from '@/services/api';
import { toast } from 'react-hot-toast';
import { Site, CreateSiteRequest } from '@/types/hosting';

interface SiteManagerProps {
  onSiteCreated?: (site: Site) => void;
}

export function SiteManager({ onSiteCreated }: SiteManagerProps) {
  const { user } = useAuth();
  const [sites, setSites] = useState<Site[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);

  const [newSite, setNewSite] = useState<CreateSiteRequest>({
    domain: '',
    custom_domain: '',
    config: {
      theme: 'default',
      custom_css: '',
      custom_js: '',
      seo_title: '',
      seo_description: '',
      seo_keywords: [],
      redirect_rules: [],
      headers: {},
      error_pages: {},
      maintenance_mode: false,
      geo_blocking: []
    },
    ssl_enabled: true,
    cdn_enabled: true
  });

  useEffect(() => {
    if (user) {
      fetchSites();
    }
  }, [user]);

  const fetchSites = async () => {
    try {
      const response = await apiClient.get('/api/v1/sites');
      setSites(response.data.sites || []);
    } catch (error) {
      console.error('Error fetching sites:', error);
      toast.error('Erro ao carregar sites');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSite = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);

    try {
      const response = await apiClient.post('/api/v1/sites', newSite);
      const site = response.data;
      
      setSites([site, ...sites]);
      setShowCreateForm(false);
      setNewSite({
        domain: '',
        custom_domain: '',
        config: {
          theme: 'default',
          custom_css: '',
          custom_js: '',
          seo_title: '',
          seo_description: '',
          seo_keywords: [],
          redirect_rules: [],
          headers: {},
          error_pages: {},
          maintenance_mode: false,
          geo_blocking: []
        },
        ssl_enabled: true,
        cdn_enabled: true
      });
      
      toast.success('Site criado com sucesso!');
      onSiteCreated?.(site);
    } catch (error: any) {
      console.error('Error creating site:', error);
      toast.error(error.response?.data?.message || 'Erro ao criar site');
    } finally {
      setCreating(false);
    }
  };

  const handleDeleteSite = async (siteId: string) => {
    if (!confirm('Tem certeza que deseja excluir este site?')) {
      return;
    }

    try {
      await apiClient.delete(`/api/v1/sites/${siteId}`);
      setSites(sites.filter(site => site.id !== siteId));
      toast.success('Site excluído com sucesso!');
    } catch (error: any) {
      console.error('Error deleting site:', error);
      toast.error(error.response?.data?.message || 'Erro ao excluir site');
    }
  };

  const getStatusBadge = (status: string) => {
    const colors = {
      active: 'bg-green-100 text-green-800',
      pending: 'bg-yellow-100 text-yellow-800',
      building: 'bg-blue-100 text-blue-800',
      error: 'bg-red-100 text-red-800',
      suspended: 'bg-gray-100 text-gray-800'
    };

    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colors[status as keyof typeof colors] || colors.pending}`}>
        {status}
      </span>
    );
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Meus Sites</h1>
        <button
          onClick={() => setShowCreateForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium"
        >
          Criar Novo Site
        </button>
      </div>

      {showCreateForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">Criar Novo Site</h2>
            <form onSubmit={handleCreateSite} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Domínio
                </label>
                <input
                  type="text"
                  value={newSite.domain}
                  onChange={(e) => setNewSite({ ...newSite, domain: e.target.value })}
                  placeholder="meu-site"
                  className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2"
                  required
                />
                <p className="text-xs text-gray-500 mt-1">
                  Será criado como: {newSite.domain || 'meu-site'}.pagemagic.app
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Domínio Personalizado (opcional)
                </label>
                <input
                  type="text"
                  value={newSite.custom_domain}
                  onChange={(e) => setNewSite({ ...newSite, custom_domain: e.target.value })}
                  placeholder="www.meusite.com"
                  className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Título SEO
                </label>
                <input
                  type="text"
                  value={newSite.config.seo_title}
                  onChange={(e) => setNewSite({ 
                    ...newSite, 
                    config: { ...newSite.config, seo_title: e.target.value }
                  })}
                  placeholder="Título do site"
                  className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Descrição SEO
                </label>
                <textarea
                  value={newSite.config.seo_description}
                  onChange={(e) => setNewSite({ 
                    ...newSite, 
                    config: { ...newSite.config, seo_description: e.target.value }
                  })}
                  placeholder="Descrição do site"
                  rows={3}
                  className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2"
                />
              </div>

              <div className="flex items-center space-x-4">
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={newSite.ssl_enabled}
                    onChange={(e) => setNewSite({ ...newSite, ssl_enabled: e.target.checked })}
                    className="rounded border-gray-300"
                  />
                  <span className="ml-2 text-sm text-gray-700">SSL Habilitado</span>
                </label>

                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={newSite.cdn_enabled}
                    onChange={(e) => setNewSite({ ...newSite, cdn_enabled: e.target.checked })}
                    className="rounded border-gray-300"
                  />
                  <span className="ml-2 text-sm text-gray-700">CDN Habilitado</span>
                </label>
              </div>

              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={creating}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
                >
                  {creating ? 'Criando...' : 'Criar Site'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {sites.map((site) => (
          <div key={site.id} className="bg-white rounded-lg shadow border p-6">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">
                  {site.domain}
                </h3>
                {site.custom_domain && (
                  <p className="text-sm text-gray-600">{site.custom_domain}</p>
                )}
              </div>
              {getStatusBadge(site.status)}
            </div>

            <div className="space-y-2 text-sm text-gray-600 mb-4">
              <p>Versão: {site.version}</p>
              {site.last_deploy && (
                <p>Último deploy: {new Date(site.last_deploy).toLocaleDateString()}</p>
              )}
              <div className="flex space-x-2">
                {site.ssl_enabled && (
                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-green-100 text-green-800">
                    SSL
                  </span>
                )}
                {site.cdn_enabled && (
                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-blue-100 text-blue-800">
                    CDN
                  </span>
                )}
              </div>
            </div>

            <div className="flex justify-between items-center">
              <a
                href={`https://${site.domain}.pagemagic.app`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:text-blue-800 text-sm font-medium"
              >
                Visualizar →
              </a>
              <div className="flex space-x-2">
                <button
                  onClick={() => {/* TODO: Navigate to editor */}}
                  className="text-gray-600 hover:text-gray-800 text-sm"
                >
                  Editar
                </button>
                <button
                  onClick={() => handleDeleteSite(site.id)}
                  className="text-red-600 hover:text-red-800 text-sm"
                >
                  Excluir
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {sites.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-400 text-6xl mb-4">🌐</div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            Nenhum site criado ainda
          </h3>
          <p className="text-gray-500 mb-4">
            Crie seu primeiro site e comece a construir sua presença online.
          </p>
          <button
            onClick={() => setShowCreateForm(true)}
            className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium"
          >
            Criar Primeiro Site
          </button>
        </div>
      )}
    </div>
  );
}
