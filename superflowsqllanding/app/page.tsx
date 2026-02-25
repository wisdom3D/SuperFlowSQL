'use client';

import { useEffect } from 'react';
import {
  Github, Monitor, Clock, Database, Settings2, PieChart, Workflow, Activity,
  BarChart3, DatabaseZap, Server, Container, RefreshCw, Layers, Presentation,
  Code2, Linkedin, ArrowRight, CheckCircle2
} from 'lucide-react';
import { DiPostgresql } from "react-icons/di";
import { SiApacheairflow } from "react-icons/si";
import { SiApachesuperset } from "react-icons/si";
import { RiTwitterXFill } from "react-icons/ri";
import { RiLinkedinLine } from "react-icons/ri";
import Scene3D from '@/components/Scene3D';

export default function Home() {
  useEffect(() => {
    // Smooth scroll for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
      anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href')!);
        if (target) {
          target.scrollIntoView({ behavior: 'smooth' });
        }
      });
    });
  }, []);

  return (
    <>
      {/* Navigation */}
      <nav className="fixed top-0 w-full z-50 glass-panel border-b border-slate-100">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            {/* Logo Mark */}
            <div className="relative w-8 h-8 flex items-center justify-center">
              {/*<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-brand-500 w-8 h-8">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>*/}
              <img src="/logo-icon.png" alt="SuperFlowSQL Logo" className="w-8 h-8" />
            </div>
            <span className="text-lg font-bold tracking-tight text-slate-900">SuperFlow<span className="text-brand-500">SQL</span></span>
          </div>

          <div className="hidden md:flex items-center gap-8 text-sm font-medium text-slate-500">
            <a href="#features" className="hover:text-brand-500 transition-colors">Features</a>
            <a href="#architecture" className="hover:text-brand-500 transition-colors">Architecture</a>
            {/* <a href="#docs" className="hover:text-brand-500 transition-colors">Docs</a> */}
            {/* <a href="#pricing" className="hover:text-brand-500 transition-colors">Pricing</a> */}
          </div>

          <div className="flex items-center gap-4">
            <a href="https://github.com/superflowsql" className="hidden md:block text-slate-400 hover:text-slate-900 transition-colors">
              <Github className="w-5 h-5" />
            </a>
            <a href="#get-started" className="bg-brand-500 hover:bg-brand-600 text-white text-sm font-semibold px-4 py-2 transition-all shadow-lg shadow-brand-500/20">
              Get Started
            </a>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative pt-32 pb-20 lg:pt-48 lg:pb-32 overflow-hidden">
        <div className="absolute inset-0 grid-bg -z-10"></div>
        <div className="absolute top-0 right-0 w-1/2 h-full bg-gradient-to-l from-brand-50 to-transparent opacity-60 -z-10"></div>

        <div className="max-w-7xl mx-auto px-6 text-center lg:text-left flex flex-col lg:flex-row items-center gap-16">
          <div className="lg:w-1/2 fade-in-up">

            <h1 className="text-5xl lg:text-7xl font-bold tracking-tight text-slate-900 leading-[1.1] mb-6">
              Data Orchestration, <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-brand-500 to-blue-400">at your fingertips.</span>
            </h1>

            <p className="text-lg text-slate-600 leading-relaxed mb-8 max-w-xl mx-auto lg:mx-0">
              Data pipelines setup, monitoring and visualization in one unified environment.
            </p>

            <div className="flex flex-col sm:flex-row items-center gap-4 justify-center lg:justify-start">
              <div className="w-full sm:w-auto px-6 py-3.5 bg-slate-900 text-white font-mono text-sm rounded-lg shadow-xl flex items-center gap-3">
                <span className="text-brand-400">$</span>
                <span>pip install superflowsql</span>
              </div>
              {/* <a href="#docs" className="w-full sm:w-auto px-8 py-3.5 bg-white border border-slate-200 text-slate-700 hover:border-slate-300 hover:bg-slate-50 font-semibold rounded-lg transition-all flex items-center justify-center gap-2 group">
                Documentation
                <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </a> */}
            </div>

            <div className="mt-10 flex items-center justify-center lg:justify-start gap-8 opacity-60 grayscale hover:grayscale-0 transition-all duration-500">
              <div className="flex items-center gap-2 font-semibold text-sm text-slate-800"><DiPostgresql className="w-4 h-4" /> PostgreSQL</div>
              <div className="flex items-center gap-2 font-semibold text-sm text-slate-800"><SiApacheairflow className="w-4 h-4" />Airflow</div>
              <div className="flex items-center gap-2 font-semibold text-sm text-slate-800"><SiApachesuperset className="w-4 h-4" /> Superset</div>
            </div>
          </div>

          {/* Hero Visual / 3D Scene */}
          <div className="lg:w-1/2 w-full fade-in-up delay-200">
            <div className="relative overflow-hidden shadow-2xl border border-slate-200">
              <div className="absolute -top-12 -right-12 w-64 h-64 bg-brand-500 rounded-full blur-[100px] opacity-20 z-0"></div>
              <div className="h-[300px] lg:h-[400px] relative z-10">
                <Scene3D />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Architecture Section */}
      <section id="architecture" className="py-24 bg-slate-50 relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center mb-20 max-w-2xl mx-auto">
            <h2 className="text-3xl lg:text-4xl font-semibold tracking-tight text-slate-900 mb-4">Unified Platform Architecture</h2>
            <p className="text-slate-600">
              A seamless integration of industry-standard tools, all orchestrated via the SuperFlowSQL CLI. We've handled the complex networking so you can focus on data.
            </p>
          </div>

          <div className="relative max-w-5xl mx-auto">
            {/* SVG Connections Background */}
            <svg className="absolute inset-0 w-full h-full pointer-events-none z-0 hidden md:block" style={{ minHeight: '400px' }}>
              {/* Gradients */}
              <defs>
                <linearGradient id="lineGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" style={{ stopColor: '#cbd5e1', stopOpacity: 1 }} />
                  <stop offset="50%" style={{ stopColor: '#198df9', stopOpacity: 1 }} />
                  <stop offset="100%" style={{ stopColor: '#cbd5e1', stopOpacity: 1 }} />
                </linearGradient>
              </defs>
              {/* Connecting Lines */}
              {/* Center to Left Top (Webserver) */}
              <path d="M512 200 L250 100" stroke="#e2e8f0" strokeWidth="2" fill="none" />
              <path d="M512 200 L250 100" stroke="#198df9" strokeWidth="2" fill="none" className="flow-line" />

              {/* Center to Left Bottom (Scheduler) */}
              <path d="M512 200 L250 300" stroke="#e2e8f0" strokeWidth="2" fill="none" />
              <path d="M512 200 L250 300" stroke="#198df9" strokeWidth="2" fill="none" className="flow-line" style={{ animationDelay: '-5s' }} />

              {/* Center to Right Top (PgAdmin) */}
              <path d="M512 200 L774 100" stroke="#e2e8f0" strokeWidth="2" fill="none" />
              <path d="M512 200 L774 100" stroke="#198df9" strokeWidth="2" fill="none" className="flow-line" style={{ animationDelay: '-10s' }} />

              {/* Center to Right Bottom (Superset) */}
              <path d="M512 200 L774 300" stroke="#e2e8f0" strokeWidth="2" fill="none" />
              <path d="M512 200 L774 300" stroke="#198df9" strokeWidth="2" fill="none" className="flow-line" style={{ animationDelay: '-15s' }} />
            </svg>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 relative z-10">

              {/* Left Column */}
              <div className="space-y-24 flex flex-col justify-center">
                {/* Airflow Webserver */}
                <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-lg hover:shadow-xl transition-shadow group">
                  <div className="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                    <Monitor className="text-blue-600" />
                  </div>
                  <h3 className="font-semibold text-slate-900 mb-2">Airflow Webserver</h3>
                  <p className="text-sm text-slate-500">Monitor pipelines, trigger DAGs, and view logs.</p>
                </div>

                {/* Airflow Scheduler */}
                <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-lg hover:shadow-xl transition-shadow group">
                  <div className="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                    <Clock className="text-blue-600" />
                  </div>
                  <h3 className="font-semibold text-slate-900 mb-2">Airflow Scheduler</h3>
                  <p className="text-sm text-slate-500">Orchestrates task execution and manages timing.</p>
                </div>
              </div>

              {/* Center Column */}
              <div className="flex items-center justify-center py-12 md:py-0">
                {/* PostgreSQL */}
                <div className="bg-white p-8 rounded-2xl border-2 border-brand-200 shadow-2xl shadow-brand-500/10 text-center w-full relative">
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-brand-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider">
                    Core
                  </div>
                  <div className="w-20 h-20 bg-brand-50 rounded-2xl flex items-center justify-center mx-auto mb-6">
                    <Database className="w-10 h-10 text-brand-600" />
                  </div>
                  <h3 className="text-xl font-bold text-slate-900 mb-2">PostgreSQL</h3>
                  <p className="text-sm text-slate-500">Central metadata store and data warehouse. The heart of the platform.</p>
                </div>
              </div>

              {/* Right Column */}
              <div className="space-y-24 flex flex-col justify-center">
                {/* PgAdmin */}
                <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-lg hover:shadow-xl transition-shadow group">
                  <div className="w-12 h-12 bg-cyan-50 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                    <Settings2 className="text-cyan-600" />
                  </div>
                  <h3 className="font-semibold text-slate-900 mb-2">PgAdmin 4</h3>
                  <p className="text-sm text-slate-500">User-friendly UI for database administration.</p>
                </div>

                {/* Superset */}
                <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-lg hover:shadow-xl transition-shadow group">
                  <div className="w-12 h-12 bg-indigo-50 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                    <PieChart className="text-indigo-600" />
                  </div>
                  <h3 className="font-semibold text-slate-900 mb-2">Apache Superset</h3>
                  <p className="text-sm text-slate-500">Business intelligence and data visualization.</p>
                </div>
              </div>

            </div>
          </div>
        </div>
      </section>

      {/* Key Features */}
      <section id="features" className="py-24 bg-white">
        <div className="max-w-7xl mx-auto px-6">
          <div className="flex flex-col md:flex-row justify-between items-end mb-16 gap-6">
            <div className="max-w-xl">
              <h2 className="text-3xl lg:text-4xl font-semibold text-slate-900 tracking-tight mb-4">Everything you need to master your data</h2>
              <p className="text-slate-600 text-lg">SuperFlowSQL abstracts the complexity of setting up a modern data stack, so you can start querying in minutes.</p>
            </div>
            <a href="#docs" className="text-brand-600 font-semibold hover:text-brand-700 flex items-center gap-1 group">
              Explore all features <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </a>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {/* Feature 1 */}
            <div className="p-8 rounded-2xl bg-slate-50 border border-slate-100 hover:border-brand-200 hover:bg-brand-50/30 transition-colors group">
              <div className="w-10 h-10 bg-white border border-slate-200 rounded-lg flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition-transform duration-300">
                <Workflow className="text-slate-700 group-hover:text-brand-600" />
              </div>
              <h3 className="text-lg font-semibold text-slate-900 mb-3">Automated Orchestration</h3>
              <p className="text-sm text-slate-600 leading-relaxed">
                Create complex DAGs in Python. Schedule jobs, manage dependencies, and handle retries automatically with Airflow.
              </p>
            </div>

            {/* Feature 2 */}
            <div className="p-8 rounded-2xl bg-slate-50 border border-slate-100 hover:border-brand-200 hover:bg-brand-50/30 transition-colors group">
              <div className="w-10 h-10 bg-white border border-slate-200 rounded-lg flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition-transform duration-300">
                <Activity className="text-slate-700 group-hover:text-brand-600" />
              </div>
              <h3 className="text-lg font-semibold text-slate-900 mb-3">Real-time Monitoring</h3>
              <p className="text-sm text-slate-600 leading-relaxed">
                Track pipeline health, catch failures instantly, and get detailed logs directly from the web interface.
              </p>
            </div>

            {/* Feature 3 */}
            <div className="p-8 rounded-2xl bg-slate-50 border border-slate-100 hover:border-brand-200 hover:bg-brand-50/30 transition-colors group">
              <div className="w-10 h-10 bg-white border border-slate-200 rounded-lg flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition-transform duration-300">
                <BarChart3 className="text-slate-700 group-hover:text-brand-600" />
              </div>
              <h3 className="text-lg font-semibold text-slate-900 mb-3">Powerful Visualization</h3>
              <p className="text-sm text-slate-600 leading-relaxed">
                Turn SQL queries into stunning dashboards with Apache Superset. Share insights across your organization.
              </p>
            </div>

            {/* Feature 4 */}
            <div className="p-8 rounded-2xl bg-slate-50 border border-slate-100 hover:border-brand-200 hover:bg-brand-50/30 transition-colors group">
              <div className="w-10 h-10 bg-white border border-slate-200 rounded-lg flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition-transform duration-300">
                <DatabaseZap className="text-slate-700 group-hover:text-brand-600" />
              </div>
              <h3 className="text-lg font-semibold text-slate-900 mb-3">Database Management</h3>
              <p className="text-sm text-slate-600 leading-relaxed">
                Full control over your PostgreSQL schemas, tables, and permissions via the integrated PgAdmin 4 interface.
              </p>
            </div>

            {/* Feature 5 */}
            <div className="p-8 rounded-2xl bg-slate-50 border border-slate-100 hover:border-brand-200 hover:bg-brand-50/30 transition-colors group">
              <div className="w-10 h-10 bg-white border border-slate-200 rounded-lg flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition-transform duration-300">
                <Server className="text-slate-700 group-hover:text-brand-600" />
              </div>
              <h3 className="text-lg font-semibold text-slate-900 mb-3">Scalable Architecture</h3>
              <p className="text-sm text-slate-600 leading-relaxed">
                Built to scale. Whether you're processing megabytes or terabytes, SuperFlowSQL grows with your data needs.
              </p>
            </div>

            {/* Feature 6 */}
            <div className="p-8 rounded-2xl bg-slate-50 border border-slate-100 hover:border-brand-200 hover:bg-brand-50/30 transition-colors group">
              <div className="w-10 h-10 bg-white border border-slate-200 rounded-lg flex items-center justify-center mb-6 shadow-sm group-hover:scale-110 transition-transform duration-300">
                <Monitor className="text-slate-700 group-hover:text-brand-600" />
              </div>
              <h3 className="text-lg font-semibold text-slate-900 mb-3">Interactive TUI</h3>
              <p className="text-sm text-slate-600 leading-relaxed">
                Manage your entire data stack without ever leaving your terminal.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How it Works */}
      <section className="py-24 border-t border-slate-100 bg-white">
        <div className="max-w-7xl mx-auto px-6">
          <h2 className="text-3xl font-semibold text-slate-900 text-center mb-16">Workflow Simplified</h2>

          <div className="relative">
            <div className="absolute left-1/2 top-0 bottom-0 w-px bg-slate-200 hidden md:block"></div>

            <div className="space-y-12">
              {/* Step 1 */}
              <div className="flex flex-col md:flex-row items-center gap-8 md:gap-16">
                <div className="md:w-1/2 text-right md:pr-8 flex flex-col items-center md:items-end">
                  <div className="w-12 h-12 bg-slate-900 text-white rounded-full flex items-center justify-center font-bold text-lg mb-4 shadow-lg">1</div>
                  <h3 className="text-xl font-semibold text-slate-900">Install & Initialize</h3>
                  <p className="text-slate-600 mt-2 max-w-sm">Install the package via pip and initialize your project structure with a single command.</p>
                </div>
                <div className="md:w-1/2 md:pl-8">
                  <div className="bg-slate-900 rounded-lg p-4 font-mono text-xs text-slate-300 shadow-xl border border-slate-800">
                    <span className="text-brand-500">$</span> pip install superflowsql<br />
                    <span className="text-brand-500">$</span> superflowsql init my-project<br />
                    <span className="text-green-500">✔ Project structure created</span><br />
                    <span className="text-green-500">✔ Environment configured</span>
                  </div>
                </div>
              </div>

              {/* Step 2 */}
              <div className="flex flex-col md:flex-row-reverse items-center gap-8 md:gap-16">
                <div className="md:w-1/2 text-left md:pl-8 flex flex-col items-center md:items-start">
                  <div className="w-12 h-12 bg-white border-2 border-slate-900 text-slate-900 rounded-full flex items-center justify-center font-bold text-lg mb-4">2</div>
                  <h3 className="text-xl font-semibold text-slate-900">Start the Stack</h3>
                  <p className="text-slate-600 mt-2 max-w-sm">Launch the entire data stack (Postgres, Airflow, Superset) using the CLI or the interactive TUI.</p>
                </div>
                <div className="md:w-1/2 md:pr-8">
                  <div className="bg-slate-900 rounded-lg p-4 font-mono text-xs text-slate-300 shadow-xl border border-slate-800">
                    <span className="text-brand-500">$</span> superflowsql start<br />
                    <span className="text-green-500">✔ PostgreSQL is ready</span><br />
                    <span className="text-green-500">✔ Airflow is running</span><br />
                    <span className="text-green-500">✔ Superset is starting...</span>
                  </div>
                </div>
              </div>

              {/* Step 3 */}
              <div className="flex flex-col md:flex-row items-center gap-8 md:gap-16">
                <div className="md:w-1/2 text-right md:pr-8 flex flex-col items-center md:items-end">
                  <div className="w-12 h-12 bg-white border-2 border-slate-900 text-slate-900 rounded-full flex items-center justify-center font-bold text-lg mb-4">3</div>
                  <h3 className="text-xl font-semibold text-slate-900">Create & Monitor</h3>
                  <p className="text-slate-600 mt-2 max-w-sm">Use the TUI to scaffold new pipelines and monitor their execution in real-time.</p>
                </div>
                <div className="md:w-1/2 md:pl-8">
                  <div className="bg-slate-900 rounded-lg p-4 font-mono text-xs text-slate-300 shadow-xl border border-slate-800">
                    <span className="text-brand-500">$</span> superflowsql create-pipeline<br />
                    <span className="text-slate-400">? Enter pipeline name:</span> <span className="text-brand-300">daily_sales</span><br />
                    <span className="text-green-500">✔ Pipeline 'daily_sales' created!</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Use Cases */}
      <section className="py-24 bg-slate-50">
        <div className="max-w-7xl mx-auto px-6">
          <h2 className="text-3xl font-semibold text-slate-900 mb-12">Built for diverse data needs</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="bg-white p-6 rounded-xl border border-slate-100 hover:shadow-lg transition-all">
              <div className="h-10 w-10 bg-orange-100 rounded-full flex items-center justify-center mb-4">
                <RefreshCw className="text-orange-600 w-5 h-5" />
              </div>
              <h4 className="font-semibold text-slate-900 mb-2">ETL Automation</h4>
              <p className="text-sm text-slate-500">Extract logs, transform schemas, and load clean data efficiently.</p>
            </div>
            <div className="bg-white p-6 rounded-xl border border-slate-100 hover:shadow-lg transition-all">
              <div className="h-10 w-10 bg-purple-100 rounded-full flex items-center justify-center mb-4">
                <Layers className="text-purple-600 w-5 h-5" />
              </div>
              <h4 className="font-semibold text-slate-900 mb-2">Data Warehousing</h4>
              <p className="text-sm text-slate-500">Centralize distributed data sources into a single source of truth.</p>
            </div>
            <div className="bg-white p-6 rounded-xl border border-slate-100 hover:shadow-lg transition-all">
              <div className="h-10 w-10 bg-pink-100 rounded-full flex items-center justify-center mb-4">
                <Presentation className="text-pink-600 w-5 h-5" />
              </div>
              <h4 className="font-semibold text-slate-900 mb-2">BI & Reporting</h4>
              <p className="text-sm text-slate-500">Generate daily reports for stakeholders without manual intervention.</p>
            </div>
            <div className="bg-white p-6 rounded-xl border border-slate-100 hover:shadow-lg transition-all">
              <div className="h-10 w-10 bg-teal-100 rounded-full flex items-center justify-center mb-4">
                <Code2 className="text-teal-600 w-5 h-5" />
              </div>
              <h4 className="font-semibold text-slate-900 mb-2">Data Engineering</h4>
              <p className="text-sm text-slate-500">Test new pipeline logic in an isolated, production-like environment.</p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-24 bg-white relative overflow-hidden">
        <div className="max-w-5xl mx-auto px-6 text-center relative z-10">
          <h2 className="text-4xl lg:text-5xl font-bold tracking-tight text-slate-900 mb-6">Ready to simplify your data orchestration?</h2>
          <p className="text-lg text-slate-600 mb-10 max-w-2xl mx-auto">Get started with the SuperFlowSQL CLI and have your full stack running in minutes.</p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <div className="w-full sm:w-auto px-8 py-4 bg-slate-900 text-white font-mono text-sm rounded-lg shadow-xl flex items-center gap-3">
              <span className="text-brand-400">$</span>
              <span>pip install superflowsql</span>
            </div>
            <a href="https://github.com/superflowsql" className="w-full sm:w-auto px-8 py-4 bg-white border border-slate-200 text-slate-700 font-semibold rounded-lg hover:bg-slate-50 transition-colors flex items-center justify-center gap-2">
              <Github className="w-4 h-4" /> Star on GitHub
            </a>
          </div>
        </div>
        {/* Decorative bg elements */}
        <div className="absolute top-1/2 left-0 -translate-y-1/2 w-64 h-64 bg-brand-100 rounded-full blur-[80px] opacity-50 -z-10"></div>
        <div className="absolute top-1/2 right-0 -translate-y-1/2 w-64 h-64 bg-blue-100 rounded-full blur-[80px] opacity-50 -z-10"></div>
      </section>

      {/* Footer */}
      <footer className="bg-slate-50 border-t border-slate-200 pt-16 pb-8">
        <div className="max-w-7xl mx-auto px-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12">
            <div className="col-span-2 md:col-span-1">
              <div className="flex items-center gap-2 mb-4">
                {/* <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-brand-500 w-6 h-6">
                  <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" strokeLinecap="round" strokeLinejoin="round"/>
                </svg> */}
                <img src="/logo-icon.png" alt="SuperFlowSQL Logo" className="w-8 h-8" />
                <span className="font-bold text-slate-900">SuperFlowSQL</span>
              </div>
              <p className="text-sm text-slate-500 mb-4">Spin up a data pipelines suite in minutes.</p>
              <div className="flex gap-4">
                {/* <a href="#" className="text-slate-400 hover:text-slate-600 transition-colors"><RiTwitterXFill className="w-5 h-5" /></a> */}
                <a href="#" className="text-slate-400 hover:text-slate-600 transition-colors"><Github className="w-5 h-5" /></a>
                {/* <a href="#" className="text-slate-400 hover:text-slate-600 transition-colors"><RiLinkedinLine className="w-5 h-5" /></a> */}
              </div>
            </div>

            <div>
              <h4 className="font-semibold text-slate-900 mb-4">Product</h4>
              <ul className="space-y-2 text-sm text-slate-500">
                <li><a href="#features" className="hover:text-brand-600 transition-colors">Features</a></li>
                <li><a href="#architecture" className="hover:text-brand-600 transition-colors">Integrations</a></li>
                {/* <li><a href="#" className="hover:text-brand-600 transition-colors">Documentation</a></li>
                <li><a href="#" className="hover:text-brand-600 transition-colors">Changelog</a></li> */}
              </ul>
            </div>

            {/* <div>
              <h4 className="font-semibold text-slate-900 mb-4">Resources</h4>
              <ul className="space-y-2 text-sm text-slate-500">
                <li><a href="#" className="hover:text-brand-600 transition-colors">Community</a></li>
                <li><a href="#" className="hover:text-brand-600 transition-colors">Help Center</a></li>
                <li><a href="#" className="hover:text-brand-600 transition-colors">API Reference</a></li>
                <li><a href="#" className="hover:text-brand-600 transition-colors">Status</a></li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-slate-900 mb-4">Legal</h4>
              <ul className="space-y-2 text-sm text-slate-500">
                <li><a href="#" className="hover:text-brand-600 transition-colors">Privacy Policy</a></li>
                <li><a href="#" className="hover:text-brand-600 transition-colors">Terms of Service</a></li>
                <li><a href="#" className="hover:text-brand-600 transition-colors">Cookie Policy</a></li>
              </ul>
            </div> */}
          </div>

          <div className="pt-8 border-t border-slate-200 flex flex-col md:flex-row justify-center items-center gap-4">
            <p className="text-sm text-slate-400">© {new Date().getFullYear()} SuperFlowSQL. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </>
  );
}
