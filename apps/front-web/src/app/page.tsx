import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Hero } from '@/components/landing/hero';
import { Features } from '@/components/landing/features';
import { Pricing } from '@/components/landing/pricing';
import { FAQ } from '@/components/landing/faq';
import { Footer } from '@/components/landing/footer';

export default function HomePage() {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Header */}
      <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 max-w-screen-2xl items-center">
          <div className="mr-4 hidden md:flex">
            <Link className="mr-6 flex items-center space-x-2" href="/">
              <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
                <span className="text-primary-foreground font-bold text-lg">P</span>
              </div>
              <span className="hidden font-bold sm:inline-block">Page Magic</span>
            </Link>
            <nav className="flex items-center space-x-6 text-sm font-medium">
              <Link
                className="transition-colors hover:text-foreground/80 text-foreground/60"
                href="/features"
              >
                Recursos
              </Link>
              <Link
                className="transition-colors hover:text-foreground/80 text-foreground/60"
                href="/pricing"
              >
                Preços
              </Link>
              <Link
                className="transition-colors hover:text-foreground/80 text-foreground/60"
                href="/docs"
              >
                Documentação
              </Link>
            </nav>
          </div>
          <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
            <nav className="flex items-center space-x-2">
              <Link href="/auth/login">
                <Button variant="ghost" size="sm">
                  Entrar
                </Button>
              </Link>
              <Link href="/auth/signup">
                <Button size="sm">
                  Começar Grátis
                </Button>
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1">
        <Hero />
        <Features />
        <Pricing />
        <FAQ />
      </main>

      <Footer />
    </div>
  );
}
