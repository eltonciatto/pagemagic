import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image, Dimensions } from 'react-native';
import { Project } from '../../types';

interface ProjectCardProps {
  project: Project;
  onPress: () => void;
  onMenu?: () => void;
}

const { width } = Dimensions.get('window');
const cardWidth = (width - 48) / 2; // 2 columns with 16px margin

export const ProjectCard: React.FC<ProjectCardProps> = ({ project, onPress, onMenu }) => {
  const getStatusColor = (status: Project['status']) => {
    switch (status) {
      case 'published':
        return '#10B981';
      case 'building':
        return '#F59E0B';
      case 'draft':
        return '#6B7280';
      default:
        return '#6B7280';
    }
  };

  const getStatusText = (status: Project['status']) => {
    switch (status) {
      case 'published':
        return 'Published';
      case 'building':
        return 'Building...';
      case 'draft':
        return 'Draft';
      default:
        return 'Unknown';
    }
  };

  return (
    <TouchableOpacity style={styles.container} onPress={onPress} activeOpacity={0.7}>
      <View style={styles.imageContainer}>
        {project.thumbnail ? (
          <Image source={{ uri: project.thumbnail }} style={styles.thumbnail} />
        ) : (
          <View style={styles.placeholderImage}>
            <Text style={styles.placeholderText}>📄</Text>
          </View>
        )}
        
        <View style={[styles.statusBadge, { backgroundColor: getStatusColor(project.status) }]}>
          <Text style={styles.statusText}>{getStatusText(project.status)}</Text>
        </View>
        
        {onMenu && (
          <TouchableOpacity style={styles.menuButton} onPress={onMenu}>
            <Text style={styles.menuIcon}>⋯</Text>
          </TouchableOpacity>
        )}
      </View>
      
      <View style={styles.content}>
        <Text style={styles.title} numberOfLines={2}>
          {project.name}
        </Text>
        
        {project.description && (
          <Text style={styles.description} numberOfLines={2}>
            {project.description}
          </Text>
        )}
        
        <View style={styles.footer}>
          <Text style={styles.date}>
            {new Date(project.updated_at).toLocaleDateString()}
          </Text>
          {project.domain && (
            <Text style={styles.domain} numberOfLines={1}>
              {project.domain}
            </Text>
          )}
        </View>
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    width: cardWidth,
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  imageContainer: {
    position: 'relative',
    height: 120,
    borderTopLeftRadius: 12,
    borderTopRightRadius: 12,
    overflow: 'hidden',
  },
  thumbnail: {
    width: '100%',
    height: '100%',
    resizeMode: 'cover',
  },
  placeholderImage: {
    width: '100%',
    height: '100%',
    backgroundColor: '#F3F4F6',
    justifyContent: 'center',
    alignItems: 'center',
  },
  placeholderText: {
    fontSize: 32,
  },
  statusBadge: {
    position: 'absolute',
    top: 8,
    left: 8,
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 6,
  },
  statusText: {
    color: '#FFFFFF',
    fontSize: 10,
    fontWeight: '600',
  },
  menuButton: {
    position: 'absolute',
    top: 8,
    right: 8,
    width: 28,
    height: 28,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    borderRadius: 14,
    justifyContent: 'center',
    alignItems: 'center',
  },
  menuIcon: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: 'bold',
  },
  content: {
    padding: 12,
  },
  title: {
    fontSize: 14,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 4,
  },
  description: {
    fontSize: 12,
    color: '#6B7280',
    marginBottom: 8,
    lineHeight: 16,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  date: {
    fontSize: 11,
    color: '#9CA3AF',
  },
  domain: {
    fontSize: 11,
    color: '#3B82F6',
    fontWeight: '500',
    flex: 1,
    textAlign: 'right',
    marginLeft: 8,
  },
});
