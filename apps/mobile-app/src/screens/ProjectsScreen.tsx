import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
  StatusBar,
  RefreshControl,
  Dimensions,
} from 'react-native';
import { ProjectCard } from '../components/ProjectCard';
import { Project } from '../types';

const { width } = Dimensions.get('window');

// Mock data - replace with actual API calls
const mockProjects: Project[] = [
  {
    id: '1',
    name: 'My Restaurant Website',
    description: 'A beautiful landing page for my restaurant business',
    status: 'published',
    domain: 'myrestaurant.com',
    thumbnail: 'https://picsum.photos/300/200?random=1',
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-20T15:30:00Z',
  },
  {
    id: '2',
    name: 'Portfolio Site',
    description: 'Personal portfolio showcasing my work',
    status: 'building',
    created_at: '2024-01-18T09:00:00Z',
    updated_at: '2024-01-18T14:00:00Z',
  },
  {
    id: '3',
    name: 'Blog Website',
    description: 'Tech blog with articles and tutorials',
    status: 'draft',
    created_at: '2024-01-10T11:00:00Z',
    updated_at: '2024-01-16T16:45:00Z',
  },
];

export default function ProjectsScreen() {
  const [projects, setProjects] = useState<Project[]>(mockProjects);
  const [refreshing, setRefreshing] = useState(false);

  const onRefresh = React.useCallback(() => {
    setRefreshing(true);
    // Simulate API call
    setTimeout(() => {
      setRefreshing(false);
    }, 2000);
  }, []);

  const handleProjectPress = (project: Project) => {
    console.log('Open project:', project.id);
    // Navigate to project editor
  };

  const handleProjectMenu = (project: Project) => {
    console.log('Show menu for project:', project.id);
    // Show action sheet with options
  };

  const handleCreateProject = () => {
    console.log('Create new project');
    // Navigate to project creation
  };

  const renderProjects = () => {
    if (projects.length === 0) {
      return (
        <View style={styles.emptyState}>
          <Text style={styles.emptyTitle}>No Projects Yet</Text>
          <Text style={styles.emptyDescription}>
            Create your first website with AI assistance
          </Text>
          <TouchableOpacity style={styles.createButton} onPress={handleCreateProject}>
            <Text style={styles.createButtonText}>Create Project</Text>
          </TouchableOpacity>
        </View>
      );
    }

    return (
      <View style={styles.projectGrid}>
        {projects.map((project) => (
          <ProjectCard
            key={project.id}
            project={project}
            onPress={() => handleProjectPress(project)}
            onMenu={() => handleProjectMenu(project)}
          />
        ))}
      </View>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle="dark-content" backgroundColor="#FFFFFF" />
      
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerTitle}>My Projects</Text>
        <TouchableOpacity style={styles.addButton} onPress={handleCreateProject}>
          <Text style={styles.addButtonText}>+</Text>
        </TouchableOpacity>
      </View>

      {/* Projects List */}
      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
        showsVerticalScrollIndicator={false}
      >
        {renderProjects()}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F9FAFB',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 16,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 1,
    borderBottomColor: '#E5E7EB',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
  },
  addButton: {
    width: 36,
    height: 36,
    backgroundColor: '#3B82F6',
    borderRadius: 18,
    justifyContent: 'center',
    alignItems: 'center',
  },
  addButtonText: {
    color: '#FFFFFF',
    fontSize: 20,
    fontWeight: 'bold',
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    padding: 16,
  },
  projectGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
  },
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 32,
    minHeight: 400,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 8,
  },
  emptyDescription: {
    fontSize: 16,
    color: '#6B7280',
    textAlign: 'center',
    marginBottom: 24,
    lineHeight: 22,
  },
  createButton: {
    backgroundColor: '#3B82F6',
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  createButtonText: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '600',
  },
});
