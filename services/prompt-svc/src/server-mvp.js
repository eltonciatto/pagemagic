// Simplified MVP test script - JavaScript version
const express = require('express');
const cors = require('cors');

const app = express();
const port = process.env.PORT || 3001;

// Basic middleware
app.use(cors());
app.use(express.json());

// Mock OpenAI service for testing
class MockSiteGenerator {
  async generateSite(request) {
    console.log('Generating site for:', request.description);
    
    // Mock response - structured JSON for testing
    const mockSite = {
      site: {
        title: `Amazing ${request.description.split(' ')[0]} Website`,
        description: `Professional website for ${request.description}`,
        sections: [
          {
            id: 'hero',
            type: 'hero',
            title: `Welcome to Your ${request.description}`,
            content: `Transform your vision into reality with our amazing ${request.description} solution.`,
            order: 1,
            cta: {
              text: 'Get Started',
              link: '#contact',
              style: 'primary'
            }
          },
          {
            id: 'features', 
            type: 'features',
            title: 'Why Choose Us',
            content: 'Discover what makes us the best choice for your needs.',
            order: 2,
            features: [
              {
                title: 'Fast Delivery',
                description: 'Quick turnaround times for all projects'
              },
              {
                title: 'Professional Quality',
                description: 'High-quality results every time'
              },
              {
                title: 'Great Support',
                description: '24/7 customer support when you need it'
              }
            ]
          },
          {
            id: 'cta',
            type: 'cta', 
            title: 'Ready to Get Started?',
            content: 'Contact us today to begin your journey.',
            order: 3,
            cta: {
              text: 'Contact Us Now',
              link: '#contact',
              style: 'primary'
            }
          }
        ],
        theme: {
          primaryColor: '#3B82F6',
          secondaryColor: '#EF4444',
          fontFamily: 'sans',
          layout: 'modern'
        },
        metadata: {
          industry: 'general',
          target_audience: 'general',
          tone: 'professional',
          language: 'en'
        }
      }
    };

    return {
      success: true,
      data: mockSite,
      metadata: {
        requestId: `gen_${Date.now()}`,
        tokensUsed: 150,
        processingTime: 500,
        model: 'mock-gpt-4'
      }
    };
  }
}

const generator = new MockSiteGenerator();

// Main generation endpoint
app.post('/v1/generate', async (req, res) => {
  try {
    console.log('Received generation request:', req.body);

    // Basic validation
    if (!req.body.description) {
      return res.status(400).json({
        success: false,
        error: 'Description is required'
      });
    }

    if (!req.body.userId) {
      return res.status(400).json({
        success: false,
        error: 'User ID is required'
      });
    }

    // Generate site
    const result = await generator.generateSite({
      description: req.body.description,
      userId: req.body.userId,
      projectId: req.body.projectId,
      preferences: req.body.preferences
    });

    res.json(result);

  } catch (error) {
    console.error('Generation error:', error);
    res.status(500).json({
      success: false,
      error: 'Internal server error',
      metadata: {
        requestId: `error_${Date.now()}`,
        tokensUsed: 0,
        processingTime: 0,
        model: 'unknown'
      }
    });
  }
});

// Health check
app.get('/v1/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: 'prompt-svc',
    timestamp: new Date().toISOString(),
    version: 'mvp-0.1.0'
  });
});

// Start server
app.listen(port, () => {
  console.log(`🚀 Prompt service (MVP) running on port ${port}`);
  console.log(`📋 Test endpoint: http://localhost:${port}/v1/generate`);
  console.log(`❤️  Health check: http://localhost:${port}/v1/health`);
});

module.exports = app;
